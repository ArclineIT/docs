require "test_helper"

class PageTest < ActiveSupport::TestCase
  setup do
    @acme = Viewer.new(sub: "1", email: "owner@acme.test", name: "Acme", roles: [])
    @other = Viewer.new(sub: "2", email: "someone@else.test", name: "Other", roles: [])
  end

  test "a public page is visible to anyone" do
    assert pages(:welcome).visible_to?(nil, nil)
    assert pages(:welcome).visible_to?(@other, nil)
  end

  test "a customer page is visible only to that customer" do
    page = pages(:acme_notes)

    assert page.visible_to?(@acme, nil)
    assert_not page.visible_to?(@other, "managed_security")
    assert_not page.visible_to?(nil, nil)
  end

  test "a plan page is visible only to customers on that plan" do
    page = pages(:monitoring_runbook)

    assert page.visible_to?(@other, "managed_security")
    assert_not page.visible_to?(@other, "Managed Security Retainer")
    assert_not page.visible_to?(@other, "other_plan")
    assert_not page.visible_to?(@other, nil)
    assert_not page.visible_to?(nil, "managed_security")
  end

  test "visibility is built from the form fields and validated" do
    assert_equal "public", Page.visibility_from(type: "public", plan: "managed_security")
    assert_equal "plan:managed_security", Page.visibility_from(type: "plan", plan: "managed_security")
    assert_equal "customer:a@b.test", Page.visibility_from(type: "customer", customer: " A@B.test ")

    page = pages(:welcome)
    %w[ plan:nope customer: customer:not-an-email secret public:x ].each do |visibility|
      page.visibility = visibility
      assert_not page.valid?, visibility
    end
  end

  test "slugs are unique, lowercase and URL-safe" do
    assert_not Page.new(title: "X", slug: "welcome").valid?
    assert_not Page.new(title: "X", slug: "has space").valid?
    assert_not Page.new(title: "X", slug: "../admin").valid?
    assert_equal "my-page", Page.new(slug: " My-Page ").slug
  end

  test "labels describe who can see a page" do
    assert_equal "Everyone signed in", pages(:welcome).visibility_label
    assert_equal "Plan: Managed Security Retainer", pages(:monitoring_runbook).visibility_label
    assert_equal "Customer: owner@acme.test", pages(:acme_notes).visibility_label
  end
end
