package main

func (s *Session) upload(args []string, options map[string]string) int {
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}

	s.Error("Not implemented: upload ")
	return -1
}
