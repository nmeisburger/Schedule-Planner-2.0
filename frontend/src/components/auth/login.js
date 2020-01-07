import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';
import qs from 'qs';
import './auth.css';

function Login(props) {

    const config = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const history = useHistory();

    function login() {
        axios.post("http://localhost:3002/signin",
            qs.stringify({
                username: username,
                password: password
            }),
            config)
            .then(res => {
                if (res.status === 200) {
                    props.authService.isAuthenticated = true;
                    props.authService.token = res.headers["access-token"]
                    history.replace("/planner");
                }
            })
            .catch(err => { console.log(err) })
    }

    function updateUsername(event) {
        setUsername(event.target.value);
    };

    function updatePassword(event) {
        setPassword(event.target.value);
    };

    return (
        <div className="auth-box" style={{ marginTop: "10%" }}>
            <div>
                <input className="auth-input" onChange={updateUsername} placeholder="Username" value={username}></input>
            </div>
            <div>
                <input className="auth-input" onChange={updatePassword} placeholder="Password" value={password} type="password"></input>
            </div>
            <div className="auth-submit" onClick={login}>
                Login
           </div>
        </div >
    );
};

export default Login;