# Signs in by walking the real callback with a stubbed Nexus client.
module SignInHelper
  FakeNexus = Struct.new(:info, :error) do
    def authorize_url(redirect_uri:, state:, code_verifier:)
      "https://nexus.test/oauth/authorize?state=#{state}&redirect_uri=#{CGI.escape(redirect_uri)}"
    end

    def userinfo(**)
      raise Nexus::Error, error if error
      info
    end
  end

  def with_nexus(info: nil, error: nil, configured: true, &block)
    stub_method(Nexus, :configured?, configured) do
      stub_method(Nexus, :client, FakeNexus.new(info, error), &block)
    end
  end

  def sign_in_as(email:, name: "Test User", roles: [])
    with_nexus(info: { "sub" => "sub-#{email}", "email" => email, "name" => name, "roles" => roles }) do
      get sign_in_path
      state = Rack::Utils.parse_query(URI.parse(response.location).query)["state"]
      get nexus_callback_path, params: { code: "abc", state: state }
    end
  end
end

ActiveSupport.on_load(:action_dispatch_integration_test) { include SignInHelper }
