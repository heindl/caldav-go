package components

import (
	"fmt"
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/icalendar/values"
	. "github.com/taviti/check"
	"net/url"
	"testing"
	"time"
)

type EventSuite struct{}

var _ = Suite(new(EventSuite))

func TestEvent(t *testing.T) { TestingT(t) }

func (s *EventSuite) TestMissingEndMarshal(c *C) {
	now := time.Now().UTC()
	event := NewEvent("test", now)
	_, err := icalendar.Marshal(event)
	c.Assert(err, ErrorMatches, "end date or duration must be set")
}

func (s *EventSuite) TestBasicWithDurationMarshal(c *C) {
	now := time.Now().UTC()
	event := NewEventWithDuration("test", now, time.Hour)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDURATION:PT1H\r\nEND:VEVENT"
	fdate := now.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, fdate, fdate))
}

func (s *EventSuite) TestBasicWithEndMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	event := NewEventWithEnd("test", now, end)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate))
}

func (s *EventSuite) TestFullEventMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	oneDay := time.Hour * 24
	oneWeek := oneDay * 7
	event := NewEventWithEnd("1:2:3", now, end)
	uri, _ := url.Parse("http://taviti.com/some/attachment.ics")
	event.Attachment = values.NewUrl(*uri)
	event.Attendees = []*values.AttendeeContact{
		values.NewAttendeeContact("Jon Azoff", "jon@taviti.com"),
		values.NewAttendeeContact("Matthew Davie", "matthew@taviti.com"),
	}
	event.Categories = values.NewCSV("vinyasa", "level 1")
	event.Comments = values.NewComments("Great class, 5 stars!", "I love this class!")
	event.ContactInfo = values.NewCSV("Send us an email!", "<jon@taviti.com>")
	event.Created = event.DateStart
	event.Description = "An all-levels class combining strength and flexibility with breath"
	ex1 := values.NewDateTime(now.Add(oneWeek))
	ex2 := values.NewDateTime(now.Add(oneWeek * 2))
	event.ExceptionDateTimes = values.NewExceptionDateTimes(ex1, ex2)
	event.Geo = values.NewGeo(37.747643, -122.445400)
	event.LastModified = event.DateStart
	event.Location = values.NewLocation("Dolores Park")
	event.Organizer = values.NewOrganizerContact("Jon Azoff", "jon@taviti.com")
	event.Priority = 1
	event.RecurrenceId = event.DateStart
	r1 := values.NewDateTime(now.Add(oneWeek + oneDay))
	r2 := values.NewDateTime(now.Add(oneWeek*2 + oneDay))
	event.RecurrenceDateTimes = values.NewRecurrenceDateTimes(r1, r2)
	event.AddRecurrenceRules(values.NewRecurrenceRule(values.WeekRecurrenceFrequency))
	uri, _ = url.Parse("matthew@taviti.com")
	event.RelatedTo = values.NewUrl(*uri)
	event.Resources = values.NewCSV("yoga mat", "towel")
	event.Sequence = 1
	event.Status = values.TentativeEventStatus
	event.Summary = "Jon's Super-Sweaty Vinyasa 1"
	event.TimeTransparency = values.OpaqueTimeTransparency
	uri, _ = url.Parse("http://student.taviti.com/san-francisco/jonathan-azoff/vinyasa-1")
	event.Url = values.NewUrl(*uri)
	enc, err := icalendar.Marshal(event)
	if err != nil {
		c.Fatal(err.Error())
	}
	tmpl := "BEGIN:VEVENT\r\nUID:1:2:3\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nCREATED:%sZ\r\n" +
		"DESCRIPTION:An all-levels class combining strength and flexibility with breath\r\n" +
		"GEO:37.747643 -122.445400\r\nLAST-MODIFIED:%sZ\r\nLOCATION:Dolores Park\r\n" +
		"ORGANIZER;CN=\"Jon Azoff\":MAILTO:jon@taviti.com\r\nPRIORITY:1\r\nSEQUENCE:1\r\nSTATUS:TENTATIVE\r\n" +
		"SUMMARY:Jon's Super-Sweaty Vinyasa 1\r\nTRANSP:OPAQUE\r\n" +
		"URL;VALUE=URI:http://student.taviti.com/san-francisco/jonathan-azoff/vinyasa-1\r\n" +
		"RECURRENCE-ID:%sZ\r\nRRULE:FREQ=WEEKLY\r\nATTACH;VALUE=URI:http://taviti.com/some/attachment.ics\r\n" +
		"ATTENDEE;CN=\"Jon Azoff\":MAILTO:jon@taviti.com\r\nATTENDEE;CN=\"Matthew Davie\":MAILTO:matthew@taviti.com\r\n" +
		"CATEGORIES:vinyasa,level 1\r\nCOMMENT:Great class, 5 stars!\r\nCOMMENT:I love this class!\r\n" +
		"CONTACT:Send us an email!,<jon@taviti.com>\r\nEXDATE:%s,%s\r\nRDATE:%s,%s\r\n" +
		"RELATED-TO;VALUE=URI:matthew@taviti.com\r\nRESOURCES:yoga mat,towel\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate, sdate, sdate, sdate, ex1, ex2, r1, r2))
}

func (s *EventSuite) TestQualifiers(c *C) {
	now := time.Now().UTC()
	event := NewEventWithDuration("test", now, time.Hour)
	c.Assert(event.IsRecurrence(), Equals, false)
	event.RecurrenceId = values.NewDateTime(now)
	c.Assert(event.IsRecurrence(), Equals, true)
	c.Assert(event.IsOverride(), Equals, false)
	event.DateStart = values.NewDateTime(now.Add(time.Hour))
	c.Assert(event.IsRecurrence(), Equals, true)
	c.Assert(event.IsOverride(), Equals, true)
}

func (s *EventSuite) TestUnmarshal(c *C) {
	raw := `
BEGIN:VEVENT
DTSTART;TZID=America/Los_Angeles:20151106T110000
DTEND;TZID=America/Los_Angeles:20151106T113000
DTSTAMP:20151117T211600Z
ORGANIZER;CN=John Boiles:mailto:john@peer.com
UID:rtmk2f1vvoprehiu2cq654991o@google.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;CN=John B
 oiles;X-NUM-GUESTS=0:mailto:john@peer.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;CN=Se
 an Chan;X-NUM-GUESTS=0:mailto:sean@peer.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;CN=St
 even Chou;X-NUM-GUESTS=0:mailto:steven@peer.com
RECURRENCE-ID;TZID=America/Los_Angeles:20151103T133000
CREATED:20150615T205546Z
DESCRIPTION:
LAST-MODIFIED:20151117T211600Z
LOCATION:
SEQUENCE:2
STATUS:CONFIRMED
SUMMARY:1:1
TRANSP:OPAQUE
X-APPLE-TRAVEL-ADVISORY-BEHAVIOR:AUTOMATIC
BEGIN:VALARM
ACTION:DISPLAY
DESCRIPTION:This is an event reminder
TRIGGER:-P0DT0H10M0S
END:VALARM
BEGIN:VALARM
ACTION:DISPLAY
DESCRIPTION:This is an event reminder
TRIGGER:-P0DT0H15M0S
END:VALARM
BEGIN:VALARM
ACTION:NONE
TRIGGER;VALUE=DATE-TIME:19760401T005545Z
X-WR-ALARMUID:BA4450C6-B914-4BAD-857D-B9BA724A034D
UID:BA4450C6-B914-4BAD-857D-B9BA724A034D
ACKNOWLEDGED:20151106T193000Z
END:VALARM
END:VEVENT
`

	e := Event{}
	err := icalendar.Unmarshal(raw, &e)
	c.Assert(err, IsNil)
	c.Assert(e.Summary, Equals, "1:1")
	c.Assert(len(e.Attendees), Equals, 3)
	c.Assert(e.Attendees[0].Entry.Name, Equals, "John Boiles")
	c.Assert(e.Attendees[0].Entry.Address, Equals, "john@peer.com")
	c.Assert(e.Attendees[1].Entry.Name, Equals, "Sean Chan")
	c.Assert(e.Attendees[1].Entry.Address, Equals, "sean@peer.com")
	c.Assert(e.Attendees[2].Entry.Name, Equals, "Steven Chou")
	c.Assert(e.Attendees[2].Entry.Address, Equals, "steven@peer.com")
}
