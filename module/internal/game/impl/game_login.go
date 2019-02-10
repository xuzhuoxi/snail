package impl

func (m *ModuleGame) Login() {
	m.state.LinkCount += 1
	notifyRemotes(m, notifyState)
}

func (m *ModuleGame) Logout() {
	m.state.LinkCount -= 1
	notifyRemotes(m, notifyState)
}
