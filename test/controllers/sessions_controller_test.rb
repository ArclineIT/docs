require "test_helper"

class SessionsControllerTest < ActionDispatch::IntegrationTest
  test "sign in sends the browser to Nexus with this site's callback" do
    with_nexus do
      get sign_in_path

      assert_response :redirect
      uri = URI.parse(response.location)
      assert_equal "nexus.test", uri.host
      assert_equal "http://www.example.com/auth/nexus/callback", Rack::Utils.parse_query(uri.query)["redirect_uri"]
    end
  end

  test "a valid callback signs the user in and returns them where they were going" do
    get client_pages_path
    assert_redirected_to sign_in_path

    sign_in_as email: "Owner@Acme.test", name: "Acme Owner"
    assert_redirected_to client_pages_path

    follow_redirect!
    assert_response :success
    assert_select ".lead", text: /Acme Owner/
  end

  test "a callback with the wrong state, or none started, is refused" do
    with_nexus(info: { "sub" => "1", "email" => "a@b.test" }) do
      get nexus_callback_path, params: { code: "abc", state: "forged" }
      assert_equal "Sign-in could not be verified. Please try again.", flash[:alert]

      get sign_in_path
      get nexus_callback_path, params: { code: "abc", state: "forged" }
      assert_equal "Sign-in could not be verified. Please try again.", flash[:alert]
    end

    get client_pages_path
    assert_redirected_to sign_in_path
  end

  test "a state can only be used once" do
    with_nexus(info: { "sub" => "1", "email" => "a@b.test" }) do
      get sign_in_path
      state = Rack::Utils.parse_query(URI.parse(response.location).query)["state"]

      get nexus_callback_path, params: { code: "abc", state: state }
      delete logout_path
      get nexus_callback_path, params: { code: "abc", state: state }
      assert_equal "Sign-in could not be verified. Please try again.", flash[:alert]
    end
  end

  test "a Nexus failure or a denied consent leaves the user signed out" do
    with_nexus(error: "Nexus answered 400: invalid_grant") do
      get sign_in_path
      state = Rack::Utils.parse_query(URI.parse(response.location).query)["state"]
      get nexus_callback_path, params: { code: "abc", state: state }
      assert_equal "Sign-in failed. Please try again.", flash[:alert]

      get nexus_callback_path, params: { error: "access_denied" }
      assert_equal "Sign-in was cancelled.", flash[:alert]
    end
  end

  test "sign in is unavailable until Nexus is configured" do
    with_nexus(configured: false) do
      get sign_in_path
      assert_redirected_to root_path
      assert_equal "Sign-in is not configured.", flash[:alert]
    end
  end

  test "logging out ends the session" do
    sign_in_as email: "owner@acme.test"
    delete logout_path

    get client_pages_path
    assert_redirected_to sign_in_path
  end
end
