# The signed-in user comes from Nexus and lives in the session cookie.
module Authentication
  extend ActiveSupport::Concern

  included do
    helper_method :current_viewer, :signed_in?
  end

  private
    def current_viewer
      @current_viewer ||= Viewer.from_session(session[:viewer])
    end

    def signed_in?
      current_viewer.present?
    end

    def require_sign_in
      return if signed_in?

      session[:return_to] = request.fullpath if request.get? || request.head?
      redirect_to sign_in_path
    end

    def require_admin
      return require_sign_in unless signed_in?

      head :forbidden unless current_viewer.admin?
    end

    def sign_in(viewer)
      return_to = session[:return_to]
      reset_session
      session[:viewer] = viewer.to_session
      return_to.presence || root_path
    end
end
