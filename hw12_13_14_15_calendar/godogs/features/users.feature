Feature: user events

  Scenario: Create username test-username@yandex.ru
    Given create user with email "test-username@yandex.ru"
    When get user by email "test-username@yandex.ru"
    Then received non empty user answer

  Scenario: Check non exist username
    Given create user with email "test-username@yandex.ru"
    When get user by email "random-email@yandex.ru"
    Then received empty user answer
