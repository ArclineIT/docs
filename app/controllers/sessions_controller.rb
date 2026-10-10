class SessionsController < ApplicationController
  # Sends the browser to Nexus.
  def new
    return redirect_to root_path, alert: "Sign-in is not configured." unless Nexus.configured?

    session[:oauth_state] = SecureRandom.urlsafe_base64(24)
    session[:oauth_verifier] = SecureRandom.urlsafe_base64(48)
    redirect_to Nexus.client.authorize_url(redirect_uri: callback_url, state: session[:oauth_state], code_verifier: session[:oauth_verifier]),
                allow_other_host: true
  end

  # Nexus sends the browser back here with a code.
  def callback
    state, verifier = session.delete(:oauth_state), session.delete(:oauth_verifier)

    if params[:error].present?
      redirect_to root_path, alert: "Sign-in was cancelled."
    elsif state.blank? || !ActiveSupport::SecurityUtils.secure_compare(state, params[:state].to_s)
      redirect_to root_path, alert: "Sign-in could not be verified. Please try again."
    else
      info = Nexus.client.userinfo(code: params[:code].to_s, redirect_uri: callback_url, code_verifier: verifier)
      redirect_to sign_in(Viewer.from_userinfo(info))
    end
  rescue Nexus::Error => e
    Rails.logger.error("nexus sign-in failed: #{e.message}")
    redirect_to root_path, alert: "Sign-in failed. Please try again."
  end

  def destroy
    reset_session
    redirect_to root_path, status: :see_other
  end

  private
    def callback_url
      "#{request.base_url}/auth/nexus/callback"
    end
end
