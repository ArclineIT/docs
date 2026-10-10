class ApplicationController < ActionController::Base
  include Authentication

  # Changes to the importmap will invalidate the etag for HTML responses
  stale_when_importmap_changes

  private
    def not_found
      raise ActionController::RoutingError, "Not Found"
    end
end
