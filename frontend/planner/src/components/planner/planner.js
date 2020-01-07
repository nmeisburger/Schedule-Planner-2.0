import React, { useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';
import qs from 'qs';
import Search from '../search/search';
import { Box, Flex } from 'rebass';
import './planner.css';

import Courses from '../../test';

function Planner(props) {
	const [courses, setCourses] = useState([]);

	const history = useHistory();

	axios.defaults.headers.common['Access-Token'] = props.authService.token;


	useEffect(() => {
		if (!props.authService.isAuthenticated) {
			history.replace("/login");
		}
		axios.get("http://localhost:3002/workspace")
			.then(res => { setCourses(res.data) })
			.catch(err => { console.log(err) })
	}, [])

	function addCourse(course) {
		axios.post("http://localhost:3002/addcourse",
			qs.stringify({
				crn: course.Crn
			}))
			.then(res => {
				if (res.status === 202) {
					setCourses(c => c.concat(course));
				}
			})
			.catch(err => {
				console.log(err);
			});
	}

	const order = ['U', 'M', 'T', 'W', 'H', 'F', 'S'];

	let days = {
		'M': [[]],
		'T': [[]],
		'W': [[]],
		'H': [[]],
		'F': [[]],
		'S': [[]],
		'U': [[]]
	};

	courses.forEach(course => {
		course.Times.forEach(time => {
			const s = (time.Start % 100) / 60 * 100 + time.Start - (time.Start % 100);
			const e = (time.End % 100) / 60 * 100 + time.End - (time.End % 100);
			time.Days.forEach(day => {
				let placed = false;
				for (let channel = 0; channel < days[day].length; channel++) {
					if (days[day][channel].length > 0 && e <= days[day][channel][0].start) {
						days[day][channel].unshift({
							"course": course.SubjCode + course.CrseNum,
							"start": s,
							"end": e
						});
						placed = true;
						break;
					}
					for (let slot = 0; slot < days[day][channel].length - 1; slot++) {
						if (s >= days[day][channel][slot].end && e <= days[day][channel][slot + 1].start) {
							days[day][channel].splice(slot, 0, {
								"course": course.SubjCode + course.CrseNum,
								"start": s,
								"end": e
							});
							placed = true;
							break;
						}
					}
					if (placed) {
						break;
					}
					if (days[day][channel].length === 0 || s >= days[day][channel][days[day][channel].length - 1].end) {
						days[day][channel].push({
							"course": course.SubjCode + course.CrseNum,
							"start": s,
							"end": e
						});
						placed = true;
						break;
					}
				}
				if (!placed) {
					days[day].push([{
						"course": course.SubjCode + course.CrseNum,
						"start": s,
						"end": e
					}]);
				}
			});
		});
	});

	function removeCourse(index) {
		axios.post("http://localhost:3002/removecourse",
			qs.stringify({
				crn: courses[index].Crn
			}))
			.then(res => {
				if (res.status === 202) {
					setCourses(c => c.filter((item, i) => i !== index))
				}
			})
			.catch(err => {
				console.log(err);
			})
	}

	return (
		<>
			<Search addCourseFunc={addCourse} />
			<Box className="workspace" width={0.75} ml='auto' mr='auto'>
				<h1 className="main-title">workspace.</h1>
				<table rules="rows">
					<thead>
						<tr style={{ padding: "10px" }}>
							<th>Remove</th>
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
							courses.map(course => (
								<tr key={course.Crn + "w"}>
									<td><button onClick={() => removeCourse(courses.indexOf(course))} className="course-remove-btn">-</button></td>
									<td style={{ width: "100px" }}>{course.Title}</td>
									<td>{course.SubjCode + " " + course.CrseNum}</td>
									<td>{course.Crn}</td>
									<td>{course.Credits}</td>
									<td>{course.Dists}</td>
									<td style={{ fontSize: "0.6em", width: "140px" }}>{course.Preq}</td>
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
			<Box className="schedule" width={0.75} ml='auto' mr='auto'>
				<h1 className='main-title'>schedule.</h1>
				<Flex flexDirection='row' justifyContent='center'>
					<p style={{ width: '30px', padding: 0, margin: 0 }}></p>
					<p className='day-title'>Sunday</p>
					<p className='day-title'>Monday</p>
					<p className='day-title'>Tuesday</p>
					<p className='day-title'>Wednesday</p>
					<p className='day-title'>Thursday</p>
					<p className='day-title'>Friday</p>
					<p className='day-title'>Saturday</p>
				</Flex>
				<Flex flexDirection='row' justifyContent='center'>
					<Flex flexDirection="column">
						<p className='schedule-time'>8:00</p>
						<p className='schedule-time'>9:00</p>
						<p className='schedule-time'>10:00</p>
						<p className='schedule-time'>11:00</p>
						<p className='schedule-time'>12:00</p>
						<p className='schedule-time'>1:00</p>
						<p className='schedule-time'>2:00</p>
						<p className='schedule-time'>3:00</p>
						<p className='schedule-time'>4:00</p>
						<p className='schedule-time'>5:00</p>
						<p className='schedule-time'>6:00</p>
						<p className='schedule-time'>7:00</p>
						<p className='schedule-time'>8:00</p>
						<p className='schedule-time'>9:00</p>
						<p className='schedule-time'>10:00</p>
					</Flex>
					{
						order.map(day => {
							const numChannels = days[day].length;
							return (
								<Flex flexDirection='row' width={100} mx={'1px'}>
									{
										days[day].map(channel => {
											if (channel.length === 0) {
												return null;
											}
											let calendar = [
												<Box width={Math.floor(100 / numChannels)}
													height={(channel[0].end - channel[0].start) / 2} mt={(channel[0].start - 700) / 2 + "px"}
													style={{ backgroundColor: '#33677A', overflow: 'hidden', fontSize: '0.7em', color: '#F2F9FF' }}>{channel[0].course}
												</Box>];
											for (let i = 1; i < channel.length; i++) {
												calendar.push(
													<Box width={Math.floor(100 / numChannels)}
														height={(channel[i].end - channel[i].start) / 2} mt={(channel[i].start - channel[i - 1].end) / 2 + "px"}
														style={{ backgroundColor: '#33677A', overflow: 'hidden', fontSize: '0.7em', color: '#F2F9FF' }}>{channel[i].course}
													</Box>
												);
											}
											return (
												<Flex flexDirection='column'>
													{calendar}
												</Flex>);
										})
									}
								</Flex>
							)
						})
					}
				</Flex>

			</Box>
		</>
	)
};

export default Planner;