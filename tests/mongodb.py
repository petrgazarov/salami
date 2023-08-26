from datetime import datetime
import motor.motor_asyncio


def create_motor_client(loop):
    return motor.motor_asyncio.AsyncIOMotorClient("localhost", 27017, io_loop=loop)


async def get_current_counter(db):
    test_results = db["test_results"]
    last_result = await test_results.find_one(sort=[("created_at", -1)])
    return (last_result["test_run_counter"] + 1) if last_result else 1


async def create_test_result(
    db,
    test_run_counter: int,
    expected: str,
    actual: str,
    full_reason: str,
    result: str,
):
    test_results = db["test_results"]
    new_result = {
        "test_run_counter": test_run_counter,
        "expected": expected,
        "actual": actual,
        "full_reason": full_reason,
        "result": result,
        "created_at": datetime.utcnow(),
    }
    await test_results.insert_one(new_result)
    return new_result
