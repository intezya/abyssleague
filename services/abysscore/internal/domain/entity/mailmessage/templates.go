package mailmessage

const linkEmailCodeMessageBodyTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Confirmation</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 20px;
            background-color: #f9f9f9;
        }
        .code {
            font-size: 24px;
            font-weight: bold;
            text-align: center;
            padding: 10px;
            margin: 20px 0;
            background-color: #eee;
            border-radius: 4px;
            letter-spacing: 5px;
        }
        .footer {
            font-size: 12px;
            color: #777;
            margin-top: 30px;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <p>Yo!</p>
        
        <p>Thank you for registering for our service. To confirm your e-mail address <strong>%s</strong>, please, use use next code:</p>
        
        <div class="code">%s</div>
        
        <p>This code is valid for %d minutes. If you have not requested this code, please ignore this email.</p>
    </div>
    
    <div class="footer">
        <p>This is an automated noreply message.</p>
        <p>&copy; %d All rights reserved.</p>
    </div>
</body>
</html>
`
