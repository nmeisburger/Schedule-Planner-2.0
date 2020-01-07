import React, { useState } from 'react';
import './auth.css';

function Register() {

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    function updateUsername(event) {
        setUsername(event.target.value);
    };

    function updatePassword(event) {
        setPassword(event.target.value);
    };

    function updateConfirmPassword(event) {
        setConfirmPassword(event.target.value);
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
            <div className="auth-submit">
                Register
           </div>
        </div>
    );
};

export default Register;