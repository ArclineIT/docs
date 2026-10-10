require "sqlite3"

# Read-only view of the billing database, used to find a customer's plan and to
# list customers for the admin form. Every method degrades to "nothing" when
# the database is not configured or cannot be read.
class BillingReader
  def self.path
    Rails.configuration.x.docs.billing_db
  end

  # Name of the customer's active (or cancelling) plan, or nil.
  def self.plan_name_for(email)
    return if email.blank?

    query(<<~SQL, [ email.downcase, *Page::ACTIVE_STATUSES ]).first&.first.presence
      SELECT subscriptions.plan_name FROM subscriptions
      JOIN users ON users.id = subscriptions.user_id
      WHERE users.email_address = ? AND subscriptions.status IN (?, ?)
      ORDER BY subscriptions.created_at DESC LIMIT 1
    SQL
  end

  # [ [ "Name <email>", email ], ... ] for a select box.
  def self.customers
    query("SELECT first_name, last_name, email_address FROM users ORDER BY email_address").map do |first, last, email|
      name = [ first, last ].compact_blank.join(" ")
      [ name.present? ? "#{name} <#{email}>" : email, email ]
    end
  end

  def self.query(sql, binds = [])
    return [] if path.nil? || !File.exist?(path)

    db = SQLite3::Database.new(path, readonly: true)
    db.execute(sql, binds)
  rescue SQLite3::Exception => e
    Rails.logger.warn("billing db unavailable: #{e.message}")
    []
  ensure
    db&.close
  end
  private_class_method :query
end
