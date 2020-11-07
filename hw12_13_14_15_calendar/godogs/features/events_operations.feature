Feature: event operations

  Scenario: Create user events
    Given delete events for user id 13
    And create event with user id 13, title "test", date from "2020-09-02T17:21:50.389825+03:00", date to "0001-01-01T02:30:17+04:00"
    When get events by user id 13
    Then receive successful save status
    And save event ID as global ID

  Scenario: Update user event
    Given exist event with global ID
    When update event with global ID and title "updated title"
    Then receive successful update status

  Scenario: Delete exist event
    Given exist event with global ID
    When delete event by global ID
    Then receive successful delete status

  Scenario: Delete non exist event
    Given non exist event with ID 0
    When delete event by ID 0
    Then receive error delete status