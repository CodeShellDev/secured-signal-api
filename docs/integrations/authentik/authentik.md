---
title: Authentik
---

# Authentik

This guide will show you how to use **Secured Signal API** as an Authenticator in [authentik](https://github.com/goauthentik/authentik).

## Setup

### 1. Stage

First you need to create the SMS Authenticator Setup Stage.

Go to `Flows and Stages > Stage > Create`.

![Stage 1](/integrations/authentik/stage_1.png)

Then you need to fill in your **API TOKEN** and your **sender number** (make sure to use the `Generic` Provider).
Point the **API URL** to your Secured Signal API `/v2/send` endpoint.

![Stage 2](/integrations/authentik/stage_2.png)

### 2. Flow

Go to `Flows and Stages > Flows > Create`.

After you have created the stage you need to use it in a setup flow.
Create one like in the screenshot below.

![Flow](/integrations/authentik/flow.png)

Note down your slug, you will need it later...

Once you've done that you will have to bind the previously created stage to the flow like so:

![Binding](/integrations/authentik/binding.png)

### 3. Webhook Mapping

Now we have to create a custom **Webhook Mapping**.

Go to `Customization > Property Mappings > Create`.
And select `Webhook Mapping`.

#### Simple

![Webhook Mapping](/integrations/authentik/mapping.png)

#### Advanced

For advanced setups or if you want to manage message content with Secured Signal API you may use this Webhook Mapping instead.

<details>
  <summary>Click to see screenshot</summary>

    ![Advanced Webhook Mapping](/integrations/authentik/advanced-mapping.png)

</details>

```python
return {
    "recipients": [device.phone_number],
    "token": f"{token}",
    "number": f"stage.from_number}"
}
```

> [!TIP]
> Take a look at authentiks [expression documentation](https://next.goauthentik.io/add-secure-apps/providers/property-mappings/expression) for all of the available variables.

Since you have decided to go the advanced way, you will have to use [**Message Templates**](../configuration/message-templates), here is an example:

```yaml
{{{ #://./message-template.yaml }}}
```

### 4. Enable SMS-Verification

To be able to use the newly created authenticator you need to enable **SMS-based Authenticators** in `default-authentication-mfa-validation`.

Go to `Flows and Stages > Stages` and edit the `default-authentication-mfa-validation` stage.

![MFA Settings](/integrations/authentik/mfa_stage.png)

Check `SMS-based Authenticators` and add your `signal-authentication-setup` stage.

## Register

After completing the Setup, you can finally go to `https://authentik.domain.com/if/flow/<your-slug>` and finish the SMS Authenticator Setup.

## Sources

- https://docs.goauthentik.io/add-secure-apps/flows-stages/stages/authenticator_sms
