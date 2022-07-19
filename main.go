package main

import (
	"fmt"
	"log"
	"reflect"
)

type Spell interface {
	// Name название заклинания
	Name() string
	// Char характеристика, на которую воздействует
	Char() string
	// Value количественное значение
	Value() int
}

// CastReceiver — если объект удовлетворяет этом интерфейсу, то заклинание применяется через него
type CastReceiver interface {
	ReceiveSpell(s Spell)
}

func CastToAll(spell Spell, objects []interface{}) {
	for _, obj := range objects {
		CastTo(spell, obj)
	}
}

func CastTo(spell Spell, object interface{}) {
	// Реализованная функция

	if receiver, ok := object.(CastReceiver); ok { // Обработка объектов, удовлетворяющих интерфейсу CastReceiver
		receiver.ReceiveSpell(spell)
		log.Printf("Spell %s casted to %s", spell.Name(), "type impl CastReceiver")
		return
	}

	val := reflect.ValueOf(object) // Получаем универсальный Value-объект

	// Если нам дали не указатель на структуру - завершаем работу
	// Дали не указатель - не сможем изменить значение
	// Дали не структуру - не сможем найти поле
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		log.Printf("Spell %s wasn't casted to %s %s", spell.Name(), "error type object", object)
		return
	}

	dependentField := val.Elem().FieldByName(spell.Char()) // Получаем ссылку на поле для воздействия заклинания

	if dependentField.IsValid() { // Если указанное поле существует
		if dependentField.CanSet() { // Если указанное поле можно изменить
			switch dependentField.Kind() { // Посмотрим тип указанного поля
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // Совместимый тип
				dependentField.SetInt(dependentField.Int() + int64(spell.Value())) // Применим заклинание
				log.Printf("Spell %s casted to %s", spell.Name(), val.Elem().Type())
				return
			default: // Несовместимый тип
				log.Printf("CAST ERROR TO %s", val.Elem().Type())
				panic("OPERATION NOT ALLOWED. WRONG FILED TYPE") // Что-то пошло не так. Поле имеет несовместимый тип
			}
		}
	}

	log.Printf("Spell %s wasn't casted to %s", spell.Name(), val.Elem().Type())
}

type spell struct {
	name string
	char string
	val  int
}

func newSpell(name string, char string, val int) Spell {
	return &spell{name: name, char: char, val: val}
}

func (s spell) Name() string {
	return s.name
}

func (s spell) Char() string {
	return s.char
}

func (s spell) Value() int {
	return s.val
}

type Player struct {
	name   string
	health int
}

func (p *Player) ReceiveSpell(s Spell) {
	if s.Char() == "Health" {
		p.health += s.Value()
	}
}

type Zombie struct {
	Health int
}

type Daemon struct {
	Health int
}

type Orc struct {
	Health int
}

type Wall struct {
	Durability int
}

func main() {

	player := &Player{
		name:   "Player_1",
		health: 100,
	}

	enemies := []interface{}{
		&Zombie{Health: 1000},
		&Zombie{Health: 1000},
		&Orc{Health: 500},
		&Orc{Health: 500},
		&Orc{Health: 500},
		&Daemon{Health: 1000},
		&Daemon{Health: 1000},
		&Wall{Durability: 100},
	}

	CastToAll(newSpell("fire", "Health", -50), append(enemies, player))
	CastToAll(newSpell("heal", "Health", 190), append(enemies, player))

	fmt.Println(player)
}
