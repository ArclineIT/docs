class Admin::PagesController < ApplicationController
  before_action :require_admin
  before_action :set_page, only: %i[ edit update destroy ]

  def index
    @pages = Page.ordered
  end

  def new
    @page = Page.new
  end

  def create
    @page = Page.new(page_params)

    if @page.save
      redirect_to admin_pages_path, notice: "#{@page.title} created."
    else
      render :new, status: :unprocessable_entity
    end
  end

  def edit
  end

  def update
    if @page.update(page_params)
      redirect_to admin_pages_path, notice: "#{@page.title} saved."
    else
      render :edit, status: :unprocessable_entity
    end
  end

  def destroy
    @page.destroy
    redirect_to admin_pages_path, notice: "#{@page.title} deleted.", status: :see_other
  end

  private
    def set_page
      @page = Page.find(params[:id])
    end

    def page_params
      form = params.expect(page: %i[ title slug description section content display_order visibility_type visibility_plan visibility_customer ])
      visibility = Page.visibility_from(type: form.delete(:visibility_type), plan: form.delete(:visibility_plan), customer: form.delete(:visibility_customer))
      form.merge(visibility: visibility)
    end
end
