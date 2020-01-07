import React from 'react';
import { Link } from 'react-router-dom';
import { Box } from 'rebass';
import './landing.css';

function Landing() {
    return (
        <div style={{ marginTop: "13%" }}>
            {/* <h1 style={{ color: "#F2F9FF", fontSize: "3em", fontStyle: "italic" }}>I can't Believe its not Banner!</h1> */}
            <Link to={"/login"} className="landing-link">
                <Box style={{ backgroundColor: "#33677A" }}>
                    <p>Login</p>
                </Box>
            </Link>
            <Link to={"/register"} className="landing-link">
                <Box style={{ backgroundColor: "#CC554F" }}>
                    <p>Register</p>
                </Box>
            </Link>
        </div>
    );
};

export default Landing;