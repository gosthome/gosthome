:: Copyright 2021 The Periph Authors. All rights reserved.
:: Use of this source code is governed under the Apache License, Version 2.0
:: that can be found in the LICENSE file.

cd %~dp0

echo "- Installing the virtualenv at ./venv/"
python -m venv venv
if %errorlevel% neq 0 exit /b %errorlevel%

echo "- Activating virtualenv"
call venv\Scripts\activate.bat
if %errorlevel% neq 0 exit /b %errorlevel%

echo "- Installing requirements"
pip install -U -r requirements.txt
if %errorlevel% neq 0 exit /b %errorlevel%

echo ""
echo "Congratulations! Everything is inside ./venv/"
echo "To access esphome, run:"
echo "  source venv/bin/activate"
