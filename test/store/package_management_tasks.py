#!/usr/bin/env python3

from scope_harness import (ScopeHarness,
                           Parameters,
                           PreviewView,
                           PreviewColumnMatcher,
                           PreviewMatcher,
                           PreviewWidgetMatcher)

def installPackage(test, categoryId, packageId):
	"""Install a package via the scope.

	This involves communicating with the DBus daemon, which will begin
	the package installation and report progress back.
	"""

	test.view.browse_department("")

	# Go to the preview of the package.
	preview = test.view.category(categoryId).result("snappy:"+packageId).tap()
	test.assertIsInstance(preview, PreviewView)

	# Begin the installation.
	widgetFound = False
	for columnWidgets in preview.widgets:
		for widget in columnWidgets:
			if widget.type == "actions":
				preview = widget.trigger("install", "")
				widgetFound = True

	test.assertTrue(widgetFound, "Expected preview to include Install action.")

	return preview

def uninstallPackage(test, categoryId, packageId):
	"""Uninstall a package via the scope.

	This involves communicating with the DBus daemon, which will begin
	the package uninstallation and report progress back.
	"""

	test.view.browse_department("")

	# Go to the preview of the package.
	preview = test.view.category(categoryId).result("snappy:"+packageId).tap()
	test.assertIsInstance(preview, PreviewView)

	# Request the uninstallation (it will need to be confirmed)
	widgetFound = False
	for columnWidgets in preview.widgets:
		for widget in columnWidgets:
			if widget.type == "actions":
				preview = widget.trigger("uninstall", "")
				widgetFound = True

	test.assertTrue(widgetFound,
	                "Expected preview to include Uninstall action.")

	return preview

def confirmUninstallPackage(test, preview):
	"""Confirm uninstallation request."""

	# Begin the uninstallation.
	widgetFound = False
	for columnWidgets in preview.widgets:
		for widget in columnWidgets:
			if widget.type == "actions":
				preview = widget.trigger("uninstall_confirm", "")
				widgetFound = True

	test.assertTrue(widgetFound,
	                "Expected preview to include Uninstall action.")

	return preview

def cancelUninstallPackage(test, preview):
	"""Cancel uninstallation request."""

	# Cancel the uninstallation.
	widgetFound = False
	for columnWidgets in preview.widgets:
		for widget in columnWidgets:
			if widget.type == "actions":
				preview = widget.trigger("uninstall_cancel", "")
				widgetFound = True

	test.assertTrue(widgetFound, "Expected preview to include Cancel action.")

	return preview


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
					"mascot": mascot,
					"attributes": [
						{"value": "✓ Installed"}
					]
				}))
			.widget(PreviewWidgetMatcher("actions")
				.type("actions")
				.data({
					"actions": [
						{
							"id": "open",
							"label": "Open"
						},
						{
							"id": "uninstall",
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

def verifyInstallingPreview(test, preview, title, subtitle, mascot,
                            description, version, size, operationNumber=1):
	"""Verify the preview for an installing package."""

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
			.widget(PreviewWidgetMatcher("install")
				.type("progress")
				.data({
					"source": {
						"dbus-name": "com.canonical.applications.WebdmPackageManager",
						"dbus-object": "/com/canonical/applications/WebdmPackageManager/operation/{}".format(operationNumber)
					}
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
							"id": "install",
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

def verifyUninstallingPreview(test, preview, title, subtitle, mascot,
                           description, version, size, operationNumber=1):
	"""Verify the preview for an uninstalling package."""

	test.assertMatchResult(PreviewColumnMatcher()
		.column(PreviewMatcher()
			.widget(PreviewWidgetMatcher("header")
				.type("header")
				.data({
					"title": title,
					"subtitle": subtitle,
					"mascot": mascot,
					"attributes": [
						{"value": "✓ Installed"}
					]
				}))
			.widget(PreviewWidgetMatcher("uninstall")
				.type("progress")
				.data({
					"source": {
						"dbus-name": "com.canonical.applications.WebdmPackageManager",
						"dbus-object": "/com/canonical/applications/WebdmPackageManager/operation/{}".format(operationNumber)
					}
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

def verifyConfirmUninstallPreview(test, preview, name):
	"""Verify the preview to confirm request to uninstall a package."""

	test.assertMatchResult(PreviewColumnMatcher()
		.column(PreviewMatcher()
			.widget(PreviewWidgetMatcher("confirm")
				.type("text")
				.data({
					"text": "Are you sure you want to uninstall {}?".format(name)
				}))
			.widget(PreviewWidgetMatcher("confirmation")
				.type("actions")
				.data({
					"actions": [
						{
							"id": "uninstall_confirm",
							"label": "Uninstall"
						},
						{
							"id": "uninstall_cancel",
							"label": "Cancel"
						}
					]
				})))
		.match(preview.widgets))
