# Machine-readable views of the public guides.
class FeedsController < ApplicationController
  def search
    expires_in 5.minutes, public: true
    render json: Guide.all.map { |guide| { title: guide.title, url: guide.path, section: guide.section_title, excerpt: guide.excerpt } }
  end

  def sitemap
    @sections = Guide.sections
    render formats: :xml
  end

  def rss
    @guides = Guide.all
    render formats: :xml
  end
end
