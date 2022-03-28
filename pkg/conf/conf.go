package conf

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

func NewSetting(dirpath string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("esTools.yaml")
	vp.AddConfigPath(dirpath)
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}
