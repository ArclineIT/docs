require "test_helper"

class ClientPagesControllerTest < ActionDispatch::IntegrationTest
  test "signed-out visitors are sent to sign in" do
    get client_page_path("welcome")
    assert_redirected_to sign_in_path
  end

  test "a customer sees public pages and their own, but not other plans or customers" do
    sign_in_as email: "owner@acme.test"

    get client_pages_path
    assert_select ".docs__card-title", text: "Welcome"
    assert_select ".docs__card-title", text: "Acme notes"
    assert_select ".docs__card-title", text: "Monitoring runbook", count: 0

    get client_page_path("acme-notes")
    assert_response :success
    assert_select ".prose", text: /Only for Acme/

    get client_page_path("monitoring-runbook")
    assert_response :forbidden
  end

  test "a plan page opens for a customer on that plan" do
    sign_in_as email: "someone@else.test"

    stub_method(BillingReader, :plan_name_for, ->(email) { "managed_security" if email == "someone@else.test" }) do
      get client_page_path("monitoring-runbook")
      assert_response :success

      get client_page_path("acme-notes")
      assert_response :forbidden
    end
  end

  test "an admin can open any page; an unknown slug is a 404" do
    sign_in_as email: "staff@arcline.test", roles: [ "admin" ]

    get client_page_path("acme-notes")
    assert_response :success
    assert_select "a", text: "Edit this page"

    get client_page_path("nope")
    assert_response :not_found
  end
end
