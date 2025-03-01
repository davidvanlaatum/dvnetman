import { expect, test } from '@playwright/test'

test('index loads', async ({ page }) => {
  await page.route('/api/v1/user/current', async (route) => {
    await route.fulfill({
      status: 200,
      json: {
        displayName: 'test',
        loggedIn: true,
      },
    })
  })
  await page.route('/api/v1/stats', async (route) => {
    await route.fulfill({
      status: 200,
      json: {
        deviceCount: 0,
        deviceTypeCount: 0,
        manufacturerCount: 0,
      },
    })
  })
  await page.goto('/ui')

  await expect(page.locator('pre')).toContainText(/deviceCount/i)
})
