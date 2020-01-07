import React, { useState, useEffect } from 'react';
import { Box } from 'rebass';
import './search.css';

const axios = require("axios")

function Search(props) {
    const [open, setOpen] = useState(false);
    const [mode, setMode] = useState(0);
    const [query, setQuery] = useState("");
    const [subj, setSubj] = useState("");
    const [dist, setDist] = useState(0);
    const [results, setResults] = useState([]);

    useEffect(() => {
        document.getElementById("search-bar").style.left = "-35%";
        document.getElementById("search-results").style.left = "100%";
    }, []);

    function handleSearch() {
        let url = "";
        let ready = true;
        switch (mode) {
            case 0:
                if (query.length > 0) {
                    url = "http://localhost:3002/search?query=" + query.trim().replace(" ", "_").toUpperCase() + "&";
                } else {
                    url = "http://localhost:3002/filter?";
                    ready = false;
                }
                break;
            case 1:
                if (query.length > 0) {
                    url = "http://localhost:3002/instructor?instructor=" + query.trim().replace(" ", "_").toLowerCase() + "&";
                } else {
                    url = "http://localhost:3002/filter?";
                    ready = false;
                }
                break
            case 2:
                if (query.length > 0) {
                    url = "http://localhost:3002/crn?crn=" + query.trim();
                } else {
                    ready = false;
                }
                break
            default:
                break;
        }
        if (subj.length === 4) {
            url += "subj=" + subj + "&";
            ready = true;
        }
        if (dist > 0) {
            url += "dist=" + dist;
            ready = true;
        }
        if (ready) {
            axios.get(url)
                .then(response => { setResults(response.data); })
                .catch(err => { console.log(err) })
        }
    };

    function updateQuery(event) {
        setQuery(event.target.value);
    };

    function updateSubj(event) {
        setSubj(event.target.value);
    };

    function searchToggle() {
        setOpen(!open)
        if (open) {
            document.getElementById("search-bar").style.left = "0%";
            document.getElementById("search-results").style.left = "35%";
        } else {
            document.getElementById("search-bar").style.left = "-35%";
            document.getElementById("search-results").style.left = "100%";
        }
    };

    function switchMode(i) {
        setMode(i);
        setQuery("");
    }

    return (
        <div className="search">
            <div className="search-open" onClick={searchToggle}>
                <img src="https://icon.now.sh/search/33677A/70" alt="SEARCH" />
            </div>
            <Box id="search-bar" className="search-bar" width={0.35}>
                <h1 className="title1">search</h1>
                <div className="switch-mode">
                    <p onClick={() => switchMode(0)} style={mode === 0 ? { backgroundColor: "#CC554F" } : { backgroundColor: "#33677A" }} > Search Courses</p>
                    <p onClick={() => switchMode(1)} style={mode === 1 ? { backgroundColor: "#CC554F" } : { backgroundColor: "#33677A" }}>Search Instructors</p>
                    <p onClick={() => switchMode(2)} style={mode === 2 ? { backgroundColor: "#CC554F" } : { backgroundColor: "#33677A" }}>Search by CRN</p>
                </div>
                <input className="input-field" value={query} onChange={updateQuery} placeholder={(() => { switch (mode) { case 0: return 'Course Title or Subject Code + Course Number'; case 1: return 'Instructor Name'; case 2: return 'CRN'; default: break; } })()} />
                {
                    mode === 0 || mode === 1
                        ? <>
                            <input className="input-field" style={{ width: "40%" }} value={subj} onChange={updateSubj} placeholder="Subject" />
                            <div className="dist-select">
                                <p onClick={() => setDist(0)} className={dist === 0 ? "dist-select-active" : "dist-select-inactive"}>Any</p>
                                <p onClick={() => setDist(1)} className={dist === 1 ? "dist-select-active" : "dist-select-inactive"}>I</p>
                                <p onClick={() => setDist(2)} className={dist === 2 ? "dist-select-active" : "dist-select-inactive"}>II</p>
                                <p onClick={() => setDist(3)} className={dist === 3 ? "dist-select-active" : "dist-select-inactive"}>III</p>
                            </div>
                        </>
                        : null
                }
                <p onClick={handleSearch} className="search-btn">Search</p>
            </Box >
            <Box id="search-results" className="search-results" width={0.65}>
                <img className="search-close" src="https://icon.now.sh/x/33677A/35" onClick={searchToggle} alt="X" />
                <h1 className="title2">results</h1>
                <Box className="search-results-box">
                    <table rules="rows">
                        <thead>
                            <tr style={{ padding: "10px" }}>
                                <th>Add</th>
                                <th colSpan="2">Course</th>
                                <th>CRN</th>
                                <th>Credits</th>
                                <th>Distribution</th>
                                <th>Prereq/Coreq</th>
                                <th>Instructors</th>
                                <th>Times</th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                results.map(course => (
                                    <tr key={course.Crn + "w"}>
                                        <td><button className="search-add-btn" onClick={() => props.addCourseFunc(course)}>+</button></td>
                                        <td>{course.Title}</td>
                                        <td>{course.SubjCode + " " + course.CrseNum}</td>
                                        <td>{course.Crn}</td>
                                        <td>{course.Credits}</td>
                                        <td>{course.Dists}</td>
                                        <td style={{ fontSize: "0.6em" }}>{course.Preq}</td>
                                        <td>{course.Instructors}</td>
                                        <td style={{ width: "120px" }}>{course.Times.map(time => (
                                            <p key={`${time.Start}-${time.End}`}>{time.Days}  {time.Start} - {time.End}</p>
                                        ))}</td>
                                    </tr>
                                ))
                            }
                        </tbody>
                    </table>
                </Box>

            </Box>
        </div >
    )
}

export default Search;