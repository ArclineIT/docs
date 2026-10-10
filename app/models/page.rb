# A page written in the admin area. Unlike guides it can be limited to the
# customers on one plan, or to a single customer.
class Page < ApplicationRecord
  PLANS = {
    "shared_starter" => "Shared Hosting Starter",
    "shared_pro" => "Shared Hosting Pro",
    "shared_business" => "Shared Hosting Business",
    "wp_starter" => "WordPress Starter",
    "wp_pro" => "WordPress Pro",
    "wp_business" => "WordPress Business",
    "vps_1" => "VPS Tier 1",
    "vps_2" => "VPS Tier 2",
    "vps_3" => "VPS Tier 3",
    "vps_4" => "VPS Tier 4"
  }.freeze
  ACTIVE_STATUSES = %w[ active cancelling ].freeze

  normalizes :title, :description, :section, with: ->(value) { value.to_s.strip }
  normalizes :slug, with: ->(slug) { slug.to_s.strip.downcase }
  normalizes :visibility, with: ->(visibility) { visibility.to_s.strip.downcase.presence || "public" }

  validates :title, presence: true
  validates :slug, presence: true, uniqueness: true, format: { with: /\A[a-z0-9]+(-[a-z0-9]+)*\z/, message: "may only use lowercase letters, digits and hyphens" }
  validate :visibility_is_well_formed

  scope :ordered, -> { order(:section, :display_order, :title) }

  # Builds the stored value from the admin form's three fields.
  def self.visibility_from(type:, plan: nil, customer: nil)
    case type
    when "plan" then "plan:#{plan}"
    when "customer" then "customer:#{customer.to_s.strip.downcase}"
    else "public"
    end
  end

  def visibility_type = visibility.split(":", 2).first
  def visibility_target = visibility.split(":", 2).second

  # +viewer+ is the signed-in user (or nil) and +plan_name+ the name of their
  # active plan from billing (or nil).
  def visible_to?(viewer, plan_name)
    case visibility_type
    when "public" then true
    when "customer" then viewer.present? && viewer.email == visibility_target
    when "plan" then viewer.present? && plan_name.present? && plan_name == PLANS[visibility_target]
    else false
    end
  end

  def visibility_label
    case visibility_type
    when "public" then "Everyone signed in"
    when "plan" then "Plan: #{PLANS.fetch(visibility_target, visibility_target)}"
    when "customer" then "Customer: #{visibility_target}"
    else visibility
    end
  end

  def html = Markdown.render(content)

  private
    def visibility_is_well_formed
      case visibility_type
      when "public" then errors.add(:visibility, "is not valid") unless visibility == "public"
      when "plan" then errors.add(:visibility, "names an unknown plan") unless PLANS.key?(visibility_target)
      when "customer" then errors.add(:visibility, "needs a customer email") unless visibility_target.to_s.match?(/\A\S+@\S+\z/)
      else errors.add(:visibility, "is not valid")
      end
    end
end
