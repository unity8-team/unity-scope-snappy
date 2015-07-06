#!/usr/bin/env python3

from scope_harness import (ScopeHarness,
                           PreviewColumnMatcher,
                           PreviewMatcher,
                           PreviewWidgetMatcher)

def verifyInstalledPreview(test, preview, title, subtitle, mascot,
                           description, version, size):
	"""Verify the preview for an installed package."""

	test.assertMatchResult(PreviewColumnMatcher()
		.column(PreviewMatcher()
			.widget(PreviewWidgetMatcher("header")
				.type("header")
				.data({
					"title": title,
					"subtitle": subtitle,
					"mascot": mascot
				}))
			.widget(PreviewWidgetMatcher("actions")
				.type("actions")
				.data({
					"actions": [
						{
							"id": 3,
							"label": "Open"
						},
						{
							"id": 2,
							"label": "Uninstall"
						}
					]
				}))
			.widget(PreviewWidgetMatcher("summary")
				.type("text")
				.data({
					"title": "Info",
					"text": description
				}))
			.widget(PreviewWidgetMatcher("updates_table")
				.type("table")
				.data({
					"title": "Updates",
					"values": [
						["Version number", version],
						["Size", size]
					]
				})))
		.match(preview.widgets))

def verifyStorePreview(test, preview, title, subtitle, mascot,
                       description, version, size):
	"""Verify the preview for an in-store package (i.e. not installed)."""

	test.assertMatchResult(PreviewColumnMatcher()
		.column(PreviewMatcher()
			.widget(PreviewWidgetMatcher("header")
				.type("header")
				.data({
					"title": title,
					"subtitle": subtitle,
					"mascot": mascot,
					"attributes": [
						{"value": "FREE"}
					]
				}))
			.widget(PreviewWidgetMatcher("actions")
				.type("actions")
				.data({
					"actions": [
						{
							"id": 1,
							"label": "Install",
						}
					]
				}))
			.widget(PreviewWidgetMatcher("summary")
				.type("text")
				.data({
					"title": "Info",
					"text": description
				}))
			.widget(PreviewWidgetMatcher("updates_table")
				.type("table")
				.data({
					"title": "Updates",
					"values": [
						["Version number", version],
						["Size", size]
					]
				})))
		.match(preview.widgets))
