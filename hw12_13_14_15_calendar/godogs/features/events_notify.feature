
Feature: check deliver messages

  Scenario: Check upcoming event
    Given delete events for user id 49
    And create event with user id 49, title "test-upcoming-event", starts in 7 days
    When get events by user id 49 after 5 seconds
    And check event notification state
    Then received true status

  Scenario: Check event that comes more than 2 week
    Given delete events for user id 53
    And create event with user id 53, title "test-non-upcoming-event", starts in 21 days
    When get events by user id 53 after 5 seconds
    And check event notification state
    Then received false status
