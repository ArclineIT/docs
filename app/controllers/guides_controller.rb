class GuidesController < ApplicationController
  def index
    @sections = Guide.sections
  end

  def section
    @section = Guide.find_section(params[:section]) or not_found
  end

  def show
    @guide = Guide.find(params[:section], params[:slug]) or not_found
  end
end
