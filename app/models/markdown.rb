# Renders the Markdown used by guides and client pages: GitHub-flavoured, with
# tables, strikethrough, task lists and heading IDs. Raw HTML is passed through,
# so only trusted authors (the repository and admins) may write it.
module Markdown
  FRONTMATTER = /\A---\n(.+?)\n---\n?/m

  def self.render(source)
    Kramdown::Document.new(source.to_s, input: "GFM", hard_wrap: false, syntax_highlighter: nil).to_html.html_safe
  end

  # Splits "key: value" frontmatter from the body. Returns [ hash, body ].
  def self.split_frontmatter(source)
    match = FRONTMATTER.match(source) or return [ {}, source ]

    fields = match[1].lines.filter_map do |line|
      key, value = line.split(":", 2)
      [ key.strip, value.strip.delete_prefix('"').delete_suffix('"') ] if value
    end.to_h
    [ fields, match.post_match ]
  end

  # Plain text of the opening of a document, for search results.
  def self.excerpt(source, length: 200)
    text = source.to_s.gsub(/```.*?```/m, " ").gsub(/^#+\s.*$/, " ").gsub(/[*_`>#|\[\]]|\(https?:[^)]*\)|^-{3,}$/, " ")
    text.squish.truncate(length, separator: " ")
  end
end
