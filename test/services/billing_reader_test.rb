require "test_helper"

class BillingReaderTest < ActiveSupport::TestCase
  setup do
    @path = Rails.root.join("tmp", "billing-test-#{SecureRandom.hex(4)}.sqlite3").to_s
    db = SQLite3::Database.new(@path)
    db.execute "CREATE TABLE users (id INTEGER PRIMARY KEY, email_address TEXT, first_name TEXT, last_name TEXT)"
    db.execute "CREATE TABLE subscriptions (id INTEGER PRIMARY KEY, user_id INTEGER, plan_name TEXT, status TEXT, created_at DATETIME)"
    db.execute "INSERT INTO users VALUES (1, 'owner@acme.test', 'Acme', 'Owner'), (2, 'gone@old.test', '', ''), (3, 'new@none.test', 'New', '')"
    db.execute "INSERT INTO subscriptions VALUES (1, 1, 'Shared Hosting Pro', 'cancelled', '2025-01-01'), (2, 1, 'VPS Tier 2', 'cancelling', '2026-01-01'), (3, 2, 'VPS Tier 1', 'cancelled', '2026-01-01')"
    db.close
  end

  teardown { File.delete(@path) if File.exist?(@path) }

  test "finds the active or cancelling plan by email" do
    with_billing(@path) do
      assert_equal "VPS Tier 2", BillingReader.plan_name_for("Owner@Acme.test")
      assert_nil BillingReader.plan_name_for("gone@old.test"), "a cancelled plan does not count"
      assert_nil BillingReader.plan_name_for("new@none.test")
      assert_nil BillingReader.plan_name_for("nobody@nowhere.test")
      assert_nil BillingReader.plan_name_for("")
    end
  end

  test "lists customers for the admin form" do
    with_billing(@path) do
      assert_equal [ [ "gone@old.test", "gone@old.test" ], [ "New <new@none.test>", "new@none.test" ], [ "Acme Owner <owner@acme.test>", "owner@acme.test" ] ],
                   BillingReader.customers
    end
  end

  test "an unset, missing or unreadable database gives nothing rather than an error" do
    with_billing(nil) { assert_nil BillingReader.plan_name_for("owner@acme.test") }
    with_billing("/nonexistent/billing.sqlite3") { assert_empty BillingReader.customers }

    File.write(@path, "this is not a database")
    with_billing(@path) { assert_nil BillingReader.plan_name_for("owner@acme.test") }
  end

  private
    def with_billing(path, &block)
      stub_method(BillingReader, :path, -> { path }, &block)
    end
end
