require "test_helper"

class AdminPagesControllerTest < ActionDispatch::IntegrationTest
  test "the admin area needs a sign-in, and then the admin role" do
    get admin_pages_path
    assert_redirected_to sign_in_path

    sign_in_as email: "owner@acme.test", roles: [ "member" ]
    get admin_pages_path
    assert_response :forbidden

    assert_no_difference -> { Page.count } do
      post admin_pages_path, params: { page: { title: "X", slug: "x", visibility_type: "public" } }
      delete admin_page_path(pages(:welcome))
    end
    assert_response :forbidden
  end

  test "the configured admin email is an admin without the role" do
    Rails.configuration.x.docs.admin_email = "blake@arcline.test"
    sign_in_as email: "Blake@Arcline.test"

    get admin_pages_path
    assert_response :success
  ensure
    Rails.configuration.x.docs.admin_email = nil
  end

  test "an admin creates a page for one plan" do
    sign_in_as email: "staff@arcline.test", roles: [ "admin" ]

    get new_admin_page_path
    assert_response :success

    assert_difference -> { Page.count } do
      post admin_pages_path, params: { page: { title: "Alert tuning", slug: "alert-tuning", content: "# Tune", display_order: 2,
                                               visibility_type: "plan", visibility_plan: "managed_security", visibility_customer: "ignored@x.test" } }
    end
    assert_redirected_to admin_pages_path
    assert_equal "plan:managed_security", Page.find_by!(slug: "alert-tuning").visibility
  end

  test "an admin edits a page down to one customer, and bad input is shown" do
    sign_in_as email: "staff@arcline.test", roles: [ "admin" ]
    page = pages(:welcome)

    get edit_admin_page_path(page)
    assert_response :success

    patch admin_page_path(page), params: { page: { title: "Welcome", slug: "welcome", visibility_type: "customer", visibility_customer: "Owner@Acme.test" } }
    assert_equal "customer:owner@acme.test", page.reload.visibility

    patch admin_page_path(page), params: { page: { title: "", slug: "acme-notes", visibility_type: "customer", visibility_customer: "" } }
    assert_response :unprocessable_entity
    assert_select ".form__errors li", text: "Title can't be blank"
    assert_select ".form__errors li", text: "Slug has already been taken"
    assert_select ".form__errors li", text: "Visibility needs a customer email"
  end

  test "an admin deletes a page" do
    sign_in_as email: "staff@arcline.test", roles: [ "admin" ]

    assert_difference -> { Page.count }, -1 do
      delete admin_page_path(pages(:acme_notes))
    end
  end
end
