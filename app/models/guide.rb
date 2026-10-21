# A public guide: one Markdown file under content/<section>/<slug>.md.
class Guide
  CONTENT_DIR = Rails.root.join("content")
  SECTIONS = {
    "security" => "Security",
    "privacy" => "Privacy & Self-Hosting",
    "reference" => "Reference"
  }.freeze

  Section = Data.define(:slug, :title, :guides) do
    def path = "/#{slug}/"
  end

  attr_reader :section, :slug, :title, :description, :order, :body

  class << self
    # Sections in display order, each with its guides in order. Read once in
    # production; re-read on every call elsewhere so edits show up.
    def sections
      return load_sections unless Rails.env.production?

      @sections ||= load_sections
    end

    def all
      sections.flat_map(&:guides)
    end

    def find_section(slug)
      sections.find { |section| section.slug == slug }
    end

    def find(section_slug, slug)
      find_section(section_slug)&.guides&.find { |guide| guide.slug == slug }
    end

    def reset!
      @sections = nil
    end

    private
      def load_sections(dir = CONTENT_DIR)
        guides = Dir.glob("*/*.md", base: dir).map { |path| new(path, File.read(File.join(dir, path))) }
        by_section = guides.group_by(&:section)

        (SECTIONS.keys | by_section.keys.sort).filter_map do |slug|
          next unless by_section[slug]
          Section.new(slug: slug, title: SECTIONS.fetch(slug, slug.titleize), guides: by_section[slug].sort_by { |g| [ g.order, g.slug ] })
        end
      end
  end

  def initialize(relative_path, source)
    @section, file = relative_path.split("/", 2)
    @slug = file.delete_suffix(".md")
    fields, @body = Markdown.split_frontmatter(source)
    @title = fields["title"].presence || @slug.titleize
    @description = fields["description"].to_s
    @order = fields["order"].to_i
  end

  def path = "/#{section}/#{slug}/"
  def section_title = SECTIONS.fetch(section, section.titleize)
  # The page already shows the title, so a leading "# Title" line is dropped.
  def html = Markdown.render(body.sub(/\A\s*# .*\n/, ""))
  def excerpt = Markdown.excerpt(body)

  def previous = neighbour(-1)
  def next = neighbour(1)

  private
    def neighbour(offset)
      guides = self.class.find_section(section).guides
      index = guides.index { |guide| guide.slug == slug } + offset
      guides[index] if index.between?(0, guides.size - 1)
    end
end
