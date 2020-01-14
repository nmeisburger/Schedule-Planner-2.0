import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';
import qs from 'qs';
import './auth.css';

function Register() {

    const config = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const history = useHistory();

    function updateUsername(event) {
        setUsername(event.target.value);
    };

    function updatePassword(event) {
        setPassword(event.target.value);
    };

    function updateConfirmPassword(event) {
        setConfirmPassword(event.target.value);
    };

    function register() {
        if (password === confirmPassword) {
            axios.post("http://localhost:3002/register",
                qs.stringify({
                    username: username,
                    password: password
                }),
                config)
                .then(res => {
                    if (res.status === 202) {
                        history.replace("/login");
                    }
                })
                .catch(err => { console.log(err) })
        }
    };

    return (
        <div className="auth-box" style={{ marginTop: "6%" }}>
            <div>
                <input className="auth-input" onChange={updateUsername} placeholder="Username" value={username}></input>
            </div>
            <div>
                <input className="auth-input" onChange={updatePassword} placeholder="Password" value={password} type="password"></input>
            </div>
            <div>
                <input className="auth-input" onChange={updateConfirmPassword} placeholder="Confirm Password" value={confirmPassword} type="password"></input>
            </div>
            <div className="auth-submit" onClick={register}>
                Register
           </div>
        </div>
    );
};

export default Register;