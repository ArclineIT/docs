require "test_helper"

class GuidesControllerTest < ActionDispatch::IntegrationTest
  test "the home page lists every section and guide" do
    get root_path

    assert_response :success
    assert_select ".docs__section h2", Guide.sections.size
    assert_select ".docs__list a", Guide.all.size
    assert_select "link[rel=canonical][href='https://docs.arcline.it/']"
  end

  test "a section page lists its guides" do
    get "/migrate"

    assert_response :success
    assert_select "h1", "Migration Guides"
    assert_select ".docs__card", Guide.find_section("migrate").guides.size
  end

  test "a guide renders with its title once, the sidebar, and prev/next links" do
    get "/migrate/from-bluehost"

    assert_response :success
    assert_select "h1", 1
    assert_select ".prose h2"
    assert_select ".docs__sidebar a[aria-current=page]", 1
    assert_select "a[rel=prev]", 1
    assert_select "a[rel=next]", 1
    assert_select "meta[name=description]"
  end

  test "every guide page renders" do
    Guide.all.each do |guide|
      get guide.path
      assert_response :success, guide.path
    end
  end

  test "unknown sections and guides are 404s" do
    get "/nope"
    assert_response :not_found

    get "/migrate/nope"
    assert_response :not_found

    get "/migrate/from-godaddy/extra"
    assert_response :not_found
  end

  test "search index, sitemap and feed cover every guide" do
    get search_index_path
    index = response.parsed_body
    assert_equal Guide.all.size, index.size
    assert_equal %w[ excerpt section title url ], index.first.keys.sort

    get sitemap_path
    assert_equal "application/xml; charset=utf-8", response.headers["Content-Type"]
    assert_equal Guide.all.size + Guide.sections.size + 1, response.body.scan("<loc>").size
    assert_includes response.body, "<loc>https://docs.arcline.it/migrate/from-godaddy/</loc>"

    get rss_path
    assert_equal Guide.all.size, response.body.scan("<item>").size
    assert_includes response.body, "<category>Privacy &amp; Self-Hosting</category>"
  end
end
