package controller

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

// ConfigSpecial holds configs that are different between normal and overtime halves
type ConfigSpecial struct {
	HalfDuration     time.Duration `yaml:"half-duration"`
	HalfTimeDuration time.Duration `yaml:"half-time-duration"`
	TimeoutDuration  time.Duration `yaml:"timeout-duration"`
	Timeouts         int           `yaml:"timeouts"`
	BreakAfter       time.Duration `yaml:"break-after"`
}

type ConfigGeometry struct {
	FieldLength                     float64 `yaml:"field-length"`
	FieldWidth                      float64 `yaml:"field-width"`
	DefenseAreaDepth                float64 `yaml:"defense-area-depth"`
	DefenseAreaWidth                float64 `yaml:"defense-area-width"`
	PlacementOffsetTouchLine        float64 `yaml:"placement-offset-touch-line"`
	PlacementOffsetGoalLine         float64 `yaml:"placement-offset-goal-line"`
	PlacementOffsetGoalLineGoalKick float64 `yaml:"placement-offset-goal-line-goal-kick"`
	PlacementOffsetDefenseArea      float64 `yaml:"placement-offset-defense-area"`
}

// ConfigGame holds configs that are valid for the whole game
type ConfigGame struct {
	YellowCardDuration        time.Duration                `yaml:"yellow-card-duration"`
	DefaultDivision           Division                     `yaml:"default-division"`
	Normal                    ConfigSpecial                `yaml:"normal"`
	Overtime                  ConfigSpecial                `yaml:"overtime"`
	TeamChoiceTimeout         time.Duration                `yaml:"team-choice-timeout"`
	DefaultGeometry           map[Division]*ConfigGeometry `yaml:"default-geometry"`
	MultipleCardStep          int                          `yaml:"multiple-card-step"`
	MultipleFoulStep          int                          `yaml:"multiple-foul-step"`
	MultiplePlacementFailures int                          `yaml:"multiple-placement-failures"`
	MaxBots                   map[Division]int             `yaml:"max-bots"`
	AutoRefProposalTimeout    time.Duration                `yaml:"auto-ref-proposal-timeout"`
}

// ConfigNetwork holds configs for network communication
type ConfigNetwork struct {
	PublishAddress string `yaml:"publish-address"`
	VisionAddress  string `yaml:"vision-address"`
}

// ConfigServer holds configs for the available server services
type ConfigServer struct {
	AutoRef ConfigServerAutoRef `yaml:"auto-ref"`
	Team    ConfigServerTeam    `yaml:"team"`
}

// ConfigServerAutoRef holds configs for the autoRef server
type ConfigServerAutoRef struct {
	Address        string `yaml:"address"`
	TrustedKeysDir string `yaml:"trusted-keys-dir"`
}

// ConfigServerTeam holds configs for the team server
type ConfigServerTeam struct {
	Address        string `yaml:"address"`
	TrustedKeysDir string `yaml:"trusted-keys-dir"`
}

// Config structure for the game controller
type Config struct {
	Network ConfigNetwork `yaml:"network"`
	Game    ConfigGame    `yaml:"game"`
	Server  ConfigServer  `yaml:"server"`
}

// LoadConfig loads a config from given file
func LoadConfig(fileName string) (config Config, err error) {

	config = DefaultConfig()

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		return
	}

	b, err := readAll(f)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		err = errors.Errorf("Could not unmarshal config file %v. %v", fileName, err)
	}

	return
}

// DefaultConfig creates a config with default values
func DefaultConfig() (c Config) {
	c.Network.PublishAddress = "224.5.23.1:10003"
	c.Network.VisionAddress = "224.5.23.2:10006"
	c.Game.YellowCardDuration = 2 * time.Minute
	c.Game.TeamChoiceTimeout = 200 * time.Millisecond
	c.Game.MultipleCardStep = 3
	c.Game.MultipleFoulStep = 3
	c.Game.MultiplePlacementFailures = 5
	c.Game.AutoRefProposalTimeout = 5 * time.Second

	c.Game.Normal.HalfDuration = 5 * time.Minute
	c.Game.Normal.HalfTimeDuration = 5 * time.Minute
	c.Game.Normal.Timeouts = 4
	c.Game.Normal.TimeoutDuration = 5 * time.Minute
	c.Game.Normal.BreakAfter = 5 * time.Minute

	c.Game.Overtime.HalfDuration = 2*time.Minute + 30*time.Second
	c.Game.Overtime.HalfTimeDuration = 2 * time.Minute
	c.Game.Overtime.Timeouts = 2
	c.Game.Overtime.TimeoutDuration = 5 * time.Minute
	c.Game.Overtime.BreakAfter = 2 * time.Minute

	c.Game.DefaultDivision = DivA

	c.Server.AutoRef.Address = ":10007"
	c.Server.AutoRef.TrustedKeysDir = "config/trusted_keys/auto_ref"
	c.Server.Team.Address = ":10008"
	c.Server.Team.TrustedKeysDir = "config/trusted_keys/team"

	c.Game.DefaultGeometry = map[Division]*ConfigGeometry{}
	c.Game.DefaultGeometry[DivA] = new(ConfigGeometry)
	c.Game.DefaultGeometry[DivA].FieldLength = 12
	c.Game.DefaultGeometry[DivA].FieldWidth = 9
	c.Game.DefaultGeometry[DivA].DefenseAreaDepth = 1.2
	c.Game.DefaultGeometry[DivA].DefenseAreaWidth = 2.4
	c.Game.DefaultGeometry[DivA].PlacementOffsetGoalLine = 0.2
	c.Game.DefaultGeometry[DivA].PlacementOffsetGoalLineGoalKick = 1.0
	c.Game.DefaultGeometry[DivA].PlacementOffsetTouchLine = 0.2
	c.Game.DefaultGeometry[DivA].PlacementOffsetDefenseArea = 1.0

	c.Game.DefaultGeometry[DivB] = new(ConfigGeometry)
	c.Game.DefaultGeometry[DivB].FieldLength = 9
	c.Game.DefaultGeometry[DivB].FieldWidth = 6
	c.Game.DefaultGeometry[DivB].DefenseAreaDepth = 1
	c.Game.DefaultGeometry[DivB].DefenseAreaWidth = 2
	c.Game.DefaultGeometry[DivB].PlacementOffsetGoalLine = 0.2
	c.Game.DefaultGeometry[DivB].PlacementOffsetGoalLineGoalKick = 1.0
	c.Game.DefaultGeometry[DivB].PlacementOffsetTouchLine = 0.2
	c.Game.DefaultGeometry[DivB].PlacementOffsetDefenseArea = 1.0

	c.Game.MaxBots = map[Division]int{}
	c.Game.MaxBots[DivA] = 8
	c.Game.MaxBots[DivB] = 6

	return
}

func readAll(f *os.File) ([]byte, error) {
	b := make([]byte, 10000)
	n, err := f.Read(b)
	if err != nil {
		return []byte{}, errors.Errorf("Can not read config files: %v", err)
	}
	if n == len(b) {
		return []byte{}, errors.New("Buffer size for reading config file is too small")
	}
	return b[:n], nil
}
