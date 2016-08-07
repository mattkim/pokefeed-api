package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

// NewPokemon constructor
func NewPokemon(db *sqlx.DB) *Pokemon {
	pokemon := &Pokemon{}
	pokemon.db = db
	pokemon.table = "pokemon"
	pokemon.hasID = true

	return pokemon
}

type PokemonRow struct {
	ID        int         `db:"id"`
	Name      string      `db:"name"`
	ImageURL  string      `db:"image_url"`
	UpdatedAt pq.NullTime `db:"updated_at"`
	CreatedAt pq.NullTime `db:"created_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

type Pokemon struct {
	Base
}

// GetById returns record by id.
func (p *Pokemon) GetAll(tx *sqlx.Tx) ([]*PokemonRow, error) {
	pokemon := []*PokemonRow{}
	query := fmt.Sprintf("SELECT * FROM %v", p.table)
	err := p.db.Select(&pokemon, query) // TODO: not sure I understand differences here.
	// TODO: Scan?
	return pokemon, err
}

// GetById returns record by id.
func (p *Pokemon) GetByID(tx *sqlx.Tx, id int) (*PokemonRow, error) {
	pokemon := &PokemonRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", p.table)
	err := p.db.Get(pokemon, query, id)

	return pokemon, err
}
