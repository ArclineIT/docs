require "test_helper"

class GuideTest < ActiveSupport::TestCase
  test "loads every guide, grouped into sections in display order" do
    sections = Guide.sections

    assert_equal %w[ security privacy reference ], sections.map(&:slug)
    assert_equal "Privacy & Self-Hosting", sections.second.title
    assert_equal Dir.glob(Rails.root.join("content/*/*.md")).size, Guide.all.size
  end

  test "guides in a section follow their order field" do
    orders = Guide.find_section("security").guides.map(&:order)
    assert_equal orders.sort, orders
  end

  test "every guide has a title and renders" do
    Guide.all.each do |guide|
      assert guide.title.present?, guide.path
      assert guide.html.present?, guide.path
      assert_no_match(/\A\s*<h1/, guide.html, "#{guide.path} repeats its title")
    end
  end

  test "finds a guide and its neighbours" do
    guide = Guide.find("security", "fail2ban")

    assert_equal "Set Up Fail2ban for SSH Brute-Force Protection", guide.title
    assert_equal "/security/fail2ban/", guide.path
    assert_nil guide.previous
    assert_equal "security", guide.next.section
    assert_equal guide.path, guide.next.previous.path
  end

  test "unknown sections and guides are nil" do
    assert_nil Guide.find_section("nope")
    assert_nil Guide.find("security", "nope")
    assert_nil Guide.find("nope", "fail2ban")
  end
end
