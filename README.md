# FicusBridge Onboarding 

# Onboarding API

**Objective**: Allow providers to onboard into FicusBridge through a third-party platform \(e.g., Acredia\). This involves tracking the provider, collecting necessary data \(government form submission\), and granting access to FicusBridge's funding optimization software.

### Key Steps:

- Consumer lands on third party plateform\(e.g Acredia\) to onboard process.
- Consumer clicks on the redirection link Fund Manager, a **handshake API** is triggered.
    - This handshake should generate a unique token that represents the provider and their session and then send it back to Acredia \(Acredia will store it\).
    - Acredia embeds the received token into the URL to allow tracking through the onboarding process.
- Then the Consumer is redirected to a page where they can download a government form.
-  On landing on the page, a form **submission status check API** is trigger which will check if the customer has already submitted the form or not.
    - if true, redirect to FicusBridge interaction page 
    - if false,  show the download/upload page.
- Consumer will fill the form and submit the form\(as a PDF\) using **upload  API  **endpoint.
- After the form is validated  access is granted to the provider to use FicusBridge's software.

### API structure:

- **POST /v1/onboard/auth**: Generate a token and initiate the onboarding process.
- **POST v1/onboard/form-upload**: Allow providers to upload their filled form.

# Interaction API\(Interface API\)

**Objective:** Once onboarding is complete, the provider interacts with FicusBridge’s software. This API will handle the redirection and interface embedding.

**Key Steps:**  

- Redirection from Acredia to FicusBridge’s software.
- When the consumer is redirected, FicusBridge should validate the token to ensure it’s legitimate and not expired.
- After token validation, Consumer approval status will be checked.
- Then, approve the redirect and load the FicusBridge white-labeled interface.

### API structure:

- **GET  /v1/interaction/check-approvalstatus**: Check the approval status of customer.

# API Flow

- **Acredia to FicusBridge**:
    - `POST /v1/onboard/auth` \(Initiate onboarding and receive token and if the user has already registered before then it will check for submission status of form if the user has already submitted the form \)
- **Customer’s Action**:
    - `POST /v1/onboard/form-upload` \(Upload filled form\)
- **Approval Status Check**
    - `GET /v1/interaction/check-approvalstatus`\(Approval status check \)
