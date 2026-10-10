# Private pages for signed-in customers, filtered by plan or customer.
class ClientPagesController < ApplicationController
  before_action :require_sign_in

  def index
    @pages = Page.ordered.select { |page| page.visible_to?(current_viewer, plan_name) }
  end

  def show
    @page = Page.find_by(slug: params[:slug]) or not_found
    head :forbidden unless @page.visible_to?(current_viewer, plan_name) || current_viewer.admin?
  end

  private
    def plan_name
      @plan_name ||= BillingReader.plan_name_for(current_viewer.email)
    end
end
