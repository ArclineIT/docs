Rails.application.configure do
  # Public address of this site, for canonical links, the sitemap and the feed.
  config.x.docs.base_url = ENV.fetch("DOCS_BASE_URL", "https://docs.arcline.it").chomp("/")
  # A signed-in user with this email can manage client pages, as can anyone
  # holding the Nexus "admin" role.
  config.x.docs.admin_email = ENV["ADMIN_EMAIL"].to_s.strip.downcase.presence
  # Billing's SQLite database, opened read-only to look up plans.
  config.x.docs.billing_db = ENV["BILLING_DB_PATH"].presence

  # Sign-in goes through Nexus (OpenID Connect).
  config.x.nexus.issuer = ENV.fetch("NEXUS_ISSUER", "https://nexus.arcline.it").chomp("/")
  config.x.nexus.client_id = ENV["NEXUS_CLIENT_ID"].presence
  config.x.nexus.client_secret = ENV["NEXUS_CLIENT_SECRET"].presence
end
