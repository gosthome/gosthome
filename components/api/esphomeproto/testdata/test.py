# Copyright 2021 The Periph Authors. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

"""Test a device using the native API component.

See https://esphome.io/components/api.html for more information.
"""

import asyncio
import argparse
import logging
import sys
import warnings
warnings.filterwarnings("ignore")

import aioesphomeapi


def main():
    parser = argparse.ArgumentParser(description=sys.modules[__name__].__doc__)
    parser.add_argument('--host', default='localhost')
    parser.add_argument('--port', default=6053, type=int)
    parser.add_argument('--password', default='')
    parser.add_argument('--noise-psk', default=None)
    parser.add_argument('--verbose', action='store_true')
    args = parser.parse_args()

    if args.verbose:
        logging.basicConfig(
            level=logging.DEBUG,
            format='%(asctime)s    %(filename)s:%(lineno)d: %(message)s')

    async def mainloop():
        """Connect to an ESPHome device and get details."""
        loop = asyncio.get_running_loop()

        # Establish connection
        api = aioesphomeapi.APIClient(args.host, args.port, args.password, noise_psk=args.noise_psk)
        try:
            await api.connect(login=True)
        except aioesphomeapi.core.APIConnectionError as e:
            print('Can\'t access %s.' % args.host, file=sys.stderr)
            print(str(e), file=sys.stderr)
            if not args.host.endswith('.local'):
                print('Maybe try %s.local?' % args.host, file=sys.stderr)
            return 1

        # Get API version of the device's firmware.
        print('API version:', api.api_version)

        # Show device details.
        device_info = await api.device_info()
        print('Device info:', device_info)

        # List all entities of the device.
        entities, services = await api.list_entities_services()
        if entities:
            print('\nEntities:')
            for entity in entities:
                print('-', entity)
        if services:
            print('\nServices:')
            for service in services:
                print('-', service)

        # Forces a CameraState too.
        # api.request_single_image()

        fut = asyncio.Future()
        expected_keys = set(c.key for c in entities)
        states = {}
        #expected_states = sum(
        #    1 for c in entities if type(c) != aioesphomeapi.CameraInfo)
        expected_states = len(entities)
        for e in entities:
            if type(e) == aioesphomeapi.ButtonInfo:
                expected_states -= 1
        if expected_states != 0:
            def cb(state):
                # Take the first state for each entity.
                if state.key in states:
                    return
                if type(state) == aioesphomeapi.CameraState:
                    # Zap the data because it's large when logged and it's
                    # potentially not deterministic.
                    state.data = b'<elided>'
                states[state.key] = state
                if len(states) == expected_states:
                    fut.set_result(True)

            api.subscribe_states(cb)
            await fut

        # Print the state ordered by key so the output is stable.
        print('\nState:')
        for _, state in sorted(states.items()):
            print('-', state)

        ls = [c for c in entities if type(c) == aioesphomeapi.LightInfo]
        if ls:
            l =ls[0]
            api.light_command(
                key=l.key,
                state=True,
                brightness = 0.5,
                rgb=(0.1, 0.2, 0.3),
                white=0.4,
                color_temperature=5000,
                transition_length=2,
                flash_length=0.5,
                effect="disco")

        #await api.subscribe_logs()
        #await api.subscribe_service_calls()
        #await api.subscribe_home_assistant_states()
        #await api.send_home_assistant_state()
        #await api.cover_command()
        #await api.fan_command()
        #await api.switch_command()
        #await api.climate_command()
        #await api.execute_service()
        #await api.request_single_stream()

        # Surprisingly this throws on Windows. I guess nobody ever called
        # it on windows or the esphome implementation doesn't close the
        # socket. I haven't found a way to trap the exception properly.
        if sys.platform != 'win32':
            await api.disconnect()
        return 0

    loop = asyncio.get_event_loop()
    return loop.run_until_complete(mainloop())


if __name__ == '__main__':
    sys.exit(main())
