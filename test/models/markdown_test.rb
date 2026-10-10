require "test_helper"

class MarkdownTest < ActiveSupport::TestCase
  test "renders tables, strikethrough, task lists, fenced code and heading ids" do
    html = Markdown.render(<<~MD)
      ## Before you start

      | Plan | Price |
      |------|-------|
      | Pro  | ~~$9~~ $7 |

      - [x] done

      ```sh
      ssh user@host
      ```
    MD

    assert_includes html, '<h2 id="before-you-start">Before you start</h2>'
    assert_includes html, "<table>"
    assert_includes html, "<del>$9</del>"
    assert_includes html, 'type="checkbox"'
    assert_includes html, "ssh user@host"
  end

  test "splits frontmatter and strips quotes" do
    fields, body = Markdown.split_frontmatter(%(---\ntitle: "Migrate: from GoDaddy"\nsection: migrate\norder: 2\n---\n\n# Body\n))

    assert_equal({ "title" => "Migrate: from GoDaddy", "section" => "migrate", "order" => "2" }, fields)
    assert_equal "\n# Body\n", body
  end

  test "a document without frontmatter is all body" do
    assert_equal [ {}, "# Just a page\n" ], Markdown.split_frontmatter("# Just a page\n")
  end

  test "excerpt is plain text without headings, code blocks or markup" do
    excerpt = Markdown.excerpt("# Title\n\nMoving **from** [GoDaddy](https://godaddy.com) is `easy`.\n\n```sh\nrm -rf /\n```\n\n## Next\n\nMore text.")

    assert_equal "Moving from GoDaddy is easy . More text.", excerpt
  end
end
