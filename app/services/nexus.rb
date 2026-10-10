# Sign-in through Nexus, Arcline's OpenID Connect provider.
module Nexus
  Error = Class.new(StandardError)

  def self.client
    config = Rails.configuration.x.nexus
    Client.new(issuer: config.issuer, client_id: config.client_id, client_secret: config.client_secret)
  end

  def self.configured?
    config = Rails.configuration.x.nexus
    config.client_id.present? && config.client_secret.present?
  end
end
