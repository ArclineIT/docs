# The signed-in user, as Nexus described them. Kept in the session cookie.
Viewer = Data.define(:sub, :email, :name, :roles) do
  def self.from_userinfo(info)
    new(sub: info["sub"].to_s, email: info["email"].to_s.strip.downcase,
        name: info["name"].presence || info["preferred_username"].presence || info["email"].to_s,
        roles: Array(info["roles"]).map(&:to_s))
  end

  def self.from_session(data)
    return if data.blank?

    new(sub: data["sub"], email: data["email"], name: data["name"], roles: Array(data["roles"]))
  end

  def to_session = to_h.stringify_keys

  def admin?
    roles.include?("admin") || (email.present? && email == Rails.configuration.x.docs.admin_email)
  end
end
