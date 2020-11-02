Feature: get events by ranges

  Scenario: Check daily events
    Given delete events for user id 17
    And create event with user id 17, title "test today", date from "2020-09-02T17:21:50.389825+03:00", date to "2020-09-02T18:21:50.389825+03:00"
    And create event with user id 17, title "test tomorrow", date from "2020-09-03T17:21:50.000000+03:00", date to "2020-09-03T18:21:50.000000+03:00"
    When get events by user id 17, date "2020-09-02T06:00:00.000000+03:00" and period "DAILY"
    Then receive successful list events with 1 events

  Scenario: Check weekly events
    Given delete events for user id 19
    And create event with user id 19, title "test today", date from "2020-09-02T17:21:50.389825+03:00", date to "2020-09-02T18:21:50.389825+03:00"
    And create event with user id 19, title "test tomorrow", date from "2020-09-03T17:21:50.000000+03:00", date to "2020-09-03T18:21:50.000000+03:00"
    And create event with user id 19, title "next week", date from "2020-09-10T17:21:50.000000+03:00", date to "2020-09-10T18:21:50.000000+03:00"
    When get events by user id 19, date "2020-09-06T12:00:00.000000+03:00" and period "WEEKLY"
    Then receive successful list events with 2 events

  Scenario: Check mothly events
    Given delete events for user id 23
    And create event with user id 23, title "prev month", date from "2019-31-12T17:21:50.389825+03:00", date to "2019-31-12T18:21:50.389825+03:00"
    And create event with user id 23, title "this month", date from "2020-01-02T17:21:50.000000+03:00", date to "2020-01-02T18:21:50.000000+03:00"
    And create event with user id 23, title "next month", date from "2020-02-02T17:21:50.000000+03:00", date to "2020-02-02T18:21:50.000000+03:00"
    When get events by user id 23, date "2020-01-15T12:00:00.000000+03:00" and period "MONTHLY"
    Then receive successful list events with 1 events
