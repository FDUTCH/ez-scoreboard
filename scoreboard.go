package scoreboard

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
)

// BuildUpdater returns new Updater.
func BuildUpdater(v any, name ...any) (Updater, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return Updater{}, fmt.Errorf("scoreboard must be struct")
	}

	if len(name) == 0 {
		return Updater{}, fmt.Errorf("name should not be empty")
	}

	// hack
	value := reflect.New(val.Type()).Elem()
	value.Set(val)

	return Updater{name: name, val: value}, nil
}

// Space represents empty line in scoreboard.
type Space struct{}

// DynamicLine is a dynamic line in the scoreboard.
// Each time Update() is called, DynamicLine will also be called.
// If an empty string is returned, the string will not be displayed on the player's side.
type DynamicLine func(p *player.Player) string

// StaticLine is a static line.
type StaticLine string

// Updater will build and send new scoreboard each time Update() is called.
type Updater struct {
	name []any
	val  reflect.Value
}

// Update updates scoreboard for the player.
func (s Updater) Update(p *player.Player) {
	t := s.val.Type()
	sc := scoreboard.New(s.name...)
	spaceIndex := 0

	for i := 0; i < s.val.NumField(); i++ {
		field := s.val.Field(i)
		if !field.CanSet() {
			continue
		}
		var str string
		switch val := field.Interface().(type) {
		case Space:
			str = fmt.Sprintf("§%d", spaceIndex%10) + strings.Repeat(" ", spaceIndex)
			spaceIndex++
		case DynamicLine:
			str = val(p)
		case StaticLine:
			str = fmt.Sprint(val)
		default:
			formating := t.Field(i).Tag.Get("scoreboard")
			if formating == "" {
				panic(fmt.Errorf("unknown type: %T", val))
			}
			str = fmt.Sprintf(formating, val)
		}
		// I know that an empty string won't be displayed on the client side anyway, but I just want this check to be here.
		if str != "" {
			_, _ = sc.WriteString(str)
		}
	}
	p.SendScoreboard(sc)
}
