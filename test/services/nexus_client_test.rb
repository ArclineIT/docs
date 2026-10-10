require "test_helper"

class NexusClientTest < ActiveSupport::TestCase
  include LocalServerHelper

  test "the authorize URL carries the client, scopes, state and a PKCE challenge" do
    url = client("https://nexus.test").authorize_url(redirect_uri: "https://docs.test/auth/nexus/callback", state: "st4te", code_verifier: "verifier")
    uri = URI.parse(url)
    query = Rack::Utils.parse_query(uri.query)

    assert_equal "https://nexus.test/oauth/authorize", "#{uri.scheme}://#{uri.host}#{uri.path}"
    assert_equal "code", query["response_type"]
    assert_equal "docs-client", query["client_id"]
    assert_equal "openid profile email roles", query["scope"]
    assert_equal "st4te", query["state"]
    assert_equal "S256", query["code_challenge_method"]
    assert_equal Base64.urlsafe_encode64(Digest::SHA256.digest("verifier"), padding: false), query["code_challenge"]
    assert_no_match "s3cret", url
  end

  test "exchanges the code, then reads UserInfo with the access token" do
    token = http_response(200, '{"access_token":"at-123","token_type":"Bearer","id_token":"x.y.z"}')
    info = http_response(200, '{"sub":"42","email":"owner@acme.test","name":"Acme Owner","roles":["member"]}')

    with_http_server(token, info) do |port, requests|
      result = client("http://127.0.0.1:#{port}").userinfo(code: "c0de", redirect_uri: "https://docs.test/cb", code_verifier: "verifier")

      assert_equal "owner@acme.test", result["email"]
      assert_match(%r{\APOST /oauth/token }, requests[0])
      body = Rack::Utils.parse_query(requests[0].split("\r\n").last)
      assert_equal({ "grant_type" => "authorization_code", "code" => "c0de", "redirect_uri" => "https://docs.test/cb",
                     "client_id" => "docs-client", "client_secret" => "s3cret", "code_verifier" => "verifier" }, body)
      assert_match(%r{\AGET /oauth/userinfo }, requests[1])
      assert_match "Authorization: Bearer at-123", requests[1]
    end
  end

  test "a rejected code, a missing token, a missing subject and an unreachable server all raise Nexus::Error" do
    with_http_server(http_response(400, '{"error":"invalid_grant","error_description":"The code has expired"}')) do |port|
      error = assert_raises(Nexus::Error) { exchange(port) }
      assert_equal "Nexus answered 400: The code has expired", error.message
    end

    with_http_server(http_response(200, "{}")) do |port|
      assert_equal "Nexus returned no access token", assert_raises(Nexus::Error) { exchange(port) }.message
    end

    with_http_server(http_response(200, '{"access_token":"at"}'), http_response(200, '{"email":"x@y.test"}')) do |port|
      assert_equal "Nexus returned no subject", assert_raises(Nexus::Error) { exchange(port) }.message
    end

    with_http_server(http_response(502, "<html>bad gateway</html>")) do |port|
      assert_match "unreadable response", assert_raises(Nexus::Error) { exchange(port) }.message
    end

    assert_match "cannot reach Nexus", assert_raises(Nexus::Error) { exchange(closed_port) }.message
  end

  private
    def client(issuer)
      Nexus::Client.new(issuer: issuer, client_id: "docs-client", client_secret: "s3cret")
    end

    def exchange(port)
      client("http://127.0.0.1:#{port}").userinfo(code: "c", redirect_uri: "https://docs.test/cb", code_verifier: "v")
    end
end
