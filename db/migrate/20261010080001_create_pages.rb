class CreatePages < ActiveRecord::Migration[8.1]
  def change
    create_table :pages do |t|
      t.string :title, null: false
      t.string :slug, null: false
      t.string :description, null: false, default: ""
      t.string :section, null: false, default: ""
      t.text :content, null: false, default: ""
      t.string :visibility, null: false, default: "public"
      t.integer :display_order, null: false, default: 0
      t.timestamps
    end
    add_index :pages, :slug, unique: true
  end
end
