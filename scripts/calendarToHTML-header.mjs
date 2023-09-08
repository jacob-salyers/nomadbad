
export function parseICS(file) {
	function expandDates(input) {
		function expandMonth(input) {
			switch (input) {
				case '01':
					return { text: 'January', ord: 0 };
				case '02':
					return { text: 'February', ord: 1 };
				case '03':
					return { text: 'March', ord: 2 };
				case '04':
					return { text: 'April', ord: 3 };
				case '05':
					return { text: 'May', ord: 4 };
				case '06':
					return { text: 'June', ord: 5 };
				case '07':
					return { text: 'July', ord:6 };
				case '08':
					return { text: 'August', ord: 7};
				case '09':
					return { text: 'September', ord: 8 };
				case '10':
					return { text: 'October', ord: 9 };
				case '11': 
					return { text: 'November', ord: 10 };
				case '12':
					return { text: 'December', ord: 11 };
			}
		}

		return {
			year: Number(input.substring(0,4)),
			month: expandMonth(input.substring(4,6)),
			day: Number(input.substring(6,8)),
			hour: Number(input.substring(9,11)),
			minute: input.substring(11,13)
		};
	}

	function expandDays(input) {
		switch (input) {
			case 'MO':
				return 'Monday';
			case 'TU':
				return 'Tuesday';
			case 'WE':
				return 'Wednesday';
			case 'TH':
				return 'Thursday';
			case 'FR':
				return 'Friday';
			case 'SA':
				return 'Saturday';
			case 'SU':
				return 'Sunday';
			default:
				throw new Error(`Unexpected day: ${input}`);
		}
	}

	const events = [{}];
	let i = 0;
	let inEvent = false;
	for (const line of file.split('\r\n')) {
		if (!inEvent) {
			inEvent = line === 'BEGIN:VEVENT';
			continue;
		}

		inEvent = !(line === 'END:VEVENT');

		if (inEvent) {
			if (i >= events.length)
				events.push({});

			const lineArr = line.split(':');
			const k = lineArr[0];
			const v = lineArr[1];
			const e = events[i];


			switch (k) {
				case 'DTSTART;TZID=America/Chicago':
				case 'DTSTART':
					e.start = expandDates(v);
					break;
				case 'DTEND;TZID=America/Chicago':
				case 'DTEND':
					e.end = expandDates(v);
					break;
				case 'RRULE':
					for (const f of v.split(';')) {
						const fArr = f.split('=');
						const k2 = fArr[0];
						const v2 = fArr[1];

						switch (k2) {
							case 'FREQ':
								e.recurring = v2.toLowerCase();
								break;
							case 'BYDAY':
								e.days = v2.split(',').map(expandDays);
								break;
							default:
								throw new Error(`Unexpected key in RRULE: ${k2}`);
						}
					}
					break;
				case 'SUMMARY':
					e.title = v;
					break;
				case 'CREATED':
					e.created = expandDates(v);
					break;
				case 'LAST-MODIFIED':
					e.updated = expandDates(v);
					break;
			}

		} else {
			if (events[i].recurring == null)
				events[i].recurring = 'never';

			i++;
		}
	}

	return events;
}

// TODO (jacob): There's a lot to consider here
export function generateHTML(events) {
	function htmlEscape(str) {
		return String(str)
			.replace(/&/g, '&amp;')
			.replace(/</g, '&lt;')
			.replace(/>/g, '&gt;')
			.replace(/"/g, '&quot;');
	}
	function upcomingEvents(events) {
		return '';
	}

	function recurringEvents(events) {
		//const months = [...new Set(events.map(e => e.start.month.text))];
		//months.sort((a,b) => b.ord - a.ord);

		let acc = `<table class="schedule">
	<tr>
		<td></td>
		<td>Sunday</td>
		<td>Monday</td>
		<td>Tuesday</td>
		<td>Wednesday</td>
		<td>Thursday</td>
		<td>Friday</td>
		<td>Saturday</td>
	</tr>`;

		const days = [
			'Sunday',
			'Monday',
			'Tuesday',
			'Wednesday',
			'Thursday',
			'Friday',
			'Saturday'
		];

		events.sort((a, b) => -(b.start.hour+b.start.minute/60) + (a.start.hour+a.start.minute/60));
		const hours = [...new Set(events.map(e => e.start.hour))];
		const minutes = [...new Set(events.map(e => e.start.minute))];
		for (const hour of hours) {
			const eventsThisHour = events.filter(e => e.start.hour === hour);
			for (const minute of minutes) {
				const eventsThisMinute = eventsThisHour.filter(e => e.start.minute === minute);
				if (eventsThisMinute.length === 0)
					continue;

				let adjustedHour, ap;
				if (hour > 12) {
					adjustedHour = hour - 12;
					ap = 'PM'; 
				} else if (hour === 12) {
					adjustedHour = hour;
					ap = 'PM';
				}else {
					adjustedHour = hour;
					ap = 'AM';
				}

				acc += `\n\t<tr>\n\t\t<td>${adjustedHour}:${minute} ${ap}</td>`;

				for (const day of days) {
					let e = eventsThisMinute.find(e => e.days.includes(day));
					acc += `\n\t\t<!--${day}--><td>${htmlEscape(e?.title ?? '')}</td>`
				}
				acc += '\n\t</tr>';
			}
		}
		acc += '\n</table>';

		return acc;		
	}

	return {
		upcoming: upcomingEvents(events.filter(el => el.recurring === 'never')),
		recurring: recurringEvents(events.filter(el => el.recurring === 'weekly'))
	};
}
