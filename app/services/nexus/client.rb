require "net/http"

module Nexus
  # The authorization-code flow with PKCE. The user's details are read from the
  # UserInfo endpoint with the access token, which was obtained directly from
  # Nexus over TLS, so the ID token is not needed.
  class Client
    SCOPES = "openid profile email roles"
    TIMEOUT = 10

    def initialize(issuer:, client_id:, client_secret:)
      @issuer = issuer
      @client_id = client_id
      @client_secret = client_secret
    end

    def self.code_challenge(verifier)
      Base64.urlsafe_encode64(Digest::SHA256.digest(verifier), padding: false)
    end

    def authorize_url(redirect_uri:, state:, code_verifier:)
      query = {
        response_type: "code", client_id: @client_id, redirect_uri: redirect_uri, scope: SCOPES, state: state,
        code_challenge: self.class.code_challenge(code_verifier), code_challenge_method: "S256"
      }
      "#{@issuer}/oauth/authorize?#{query.to_query}"
    end

    # Trades the code for an access token and returns the UserInfo claims.
    def userinfo(code:, redirect_uri:, code_verifier:)
      token = post("/oauth/token", grant_type: "authorization_code", code: code, redirect_uri: redirect_uri,
                   client_id: @client_id, client_secret: @client_secret, code_verifier: code_verifier)
      access_token = token["access_token"].presence or raise Error, "Nexus returned no access token"

      info = get("/oauth/userinfo", "Authorization" => "Bearer #{access_token}")
      raise Error, "Nexus returned no subject" if info["sub"].blank?
      info
    end

    private
      def post(path, params)
        request = Net::HTTP::Post.new(uri(path), "Accept" => "application/json")
        request.set_form_data(params)
        perform(request)
      end

      def get(path, headers)
        perform(Net::HTTP::Get.new(uri(path), headers.merge("Accept" => "application/json")))
      end

      def uri(path)
        URI.parse("#{@issuer}#{path}")
      end

      def perform(request)
        http = Net::HTTP.new(request.uri.host, request.uri.port)
        http.use_ssl = request.uri.scheme == "https"
        http.open_timeout = http.read_timeout = TIMEOUT

        response = http.request(request)
        body = JSON.parse(response.body.presence || "{}")
        return body if response.is_a?(Net::HTTPSuccess)

        raise Error, "Nexus answered #{response.code}: #{body['error_description'] || body['error'] || response.message}"
      rescue JSON::ParserError
        raise Error, "Nexus sent an unreadable response (HTTP #{response.code})"
      rescue SystemCallError, SocketError, IOError, Timeout::Error, OpenSSL::SSL::SSLError => e
        raise Error, "cannot reach Nexus: #{e.message}"
      end
  end
end
