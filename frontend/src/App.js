import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import Landing from './components/landing/landing';
import Login from './components/auth/login';
import Register from './components/auth/register';
import Planner from './components/planner/planner';
import './App.css';

function App() {

	const authService = {
		authenticated: false,
		token: ""
	}

	return (
		<Router>
			<Switch>
				<Route path="/login">
					<Login authService={authService} />
				</Route>
				<Route path="/register">
					<Register />
				</Route>
				<Route path="/planner">
					<Planner authService={authService} />
				</Route>
				<Route path="/">
					<Landing />
				</Route>
			</Switch>
		</Router>
	);
}

export default App;
