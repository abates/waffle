package waffle


type WaffleError string

func (we WaffleError) Error() string { return string(we) }
