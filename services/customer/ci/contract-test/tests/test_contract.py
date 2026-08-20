from conftest import schema


@schema.parametrize()
def test_api_conforms_to_schema(case):
    case.call_and_validate()
